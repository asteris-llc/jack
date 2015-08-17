package jack

import (
	"bufio"
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"io"
	"os/exec"
	"sync"
)

var (
	ErrStopped = errors.New("client stopped before response received")
)

type Client struct {
	// goroutine state
	cancel  func()
	context context.Context

	// req/rep state
	lock     sync.RWMutex
	messages chan *Message
	waiting  map[string]chan *Message

	// command state
	path string
	cmd  *exec.Cmd
	in   io.WriteCloser
	out  io.ReadCloser
}

func NewClient(path string) *Client {
	context, cancel := context.WithCancel(context.Background())
	return &Client{
		context:  context,
		cancel:   cancel,
		lock:     sync.RWMutex{},
		path:     path,
		messages: make(chan *Message, 1),
		waiting:  make(map[string]chan *Message),
	}
}

func (c *Client) Start() error {
	// start command
	c.cmd = exec.Command(c.path)

	inStream, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	c.in = inStream

	outStream, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.out = outStream

	err = c.cmd.Start()
	if err != nil {
		return err
	}

	// start internal state goroutines
	go c.acceptRequests()
	go c.dispatchResponses()

	return nil
}

func (c *Client) acceptRequests() {
	newline := []byte("\n")
	for {
		select {
		case message := <-c.messages:
			blob, err := json.Marshal(message)
			if err != nil {
				panic(err) // TODO: more graceful failure
			}
			c.in.Write(blob)
			c.in.Write(newline)
		case <-c.context.Done():
			c.in.Close()
			return
		}
	}
}

func (c *Client) dispatchResponses() {
	messages := make(chan *Message, 1)
	go func(out chan *Message, src io.Reader) {
		scanner := bufio.NewScanner(c.out)
		for scanner.Scan() {
			message := new(Message)
			err := json.Unmarshal(scanner.Bytes(), message)
			if err != nil {
				panic(err) // TODO: more graceful failure
			}
			out <- message
		}
	}(messages, c.out)

	for {
		select {
		case message := <-messages:
			c.lock.RLock()
			client, ok := c.waiting[message.ID]
			if ok {
				c.lock.RUnlock()
				c.lock.Lock()
				delete(c.waiting, message.ID)
				c.lock.Unlock()
				client <- message
			} else {
				c.lock.RUnlock()
			}
		case <-c.context.Done():
			c.out.Close()
			return
		}
	}
}

func (c *Client) Stop() {
	c.cancel()
	c.cmd.Wait()
}

func (c *Client) Call(method string, args ...interface{}) (interface{}, error) {
	message := &Message{
		ID:      uuid.NewRandom().String(),
		Method:  method,
		Payload: args,
	}
	c.messages <- message

	// put our listener on the queue and listen for it
	c.lock.Lock()
	results := make(chan *Message, 1)
	c.waiting[message.ID] = results
	c.lock.Unlock()

	select {
	case result := <-results:
		return result.Payload, nil
	case <-c.context.Done():
		return nil, ErrStopped
	}
}
