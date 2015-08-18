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
	lock     sync.Mutex
	messages chan *Request
	waiting  map[string]chan *Response

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
		lock:     sync.Mutex{},
		path:     path,
		messages: make(chan *Request, 1),
		waiting:  make(map[string]chan *Response),
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
	messages := make(chan *Response, 1)

	go func(out chan *Response, src io.Reader) {
		scanner := bufio.NewScanner(c.out)

		for scanner.Scan() {
			response := new(Response)
			err := json.Unmarshal(scanner.Bytes(), response)
			if err != nil {
				panic(err) // TODO: more graceful failure
			}
			out <- response
		}
	}(messages, c.out)

	for {
		select {
		case message := <-messages:
			c.lock.Lock()
			client, ok := c.waiting[message.ID]
			if ok {
				delete(c.waiting, message.ID)
				client <- message
			}
			c.lock.Unlock()

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
	message := &Request{
		ID:     uuid.NewRandom().String(),
		Method: method,
		Args:   args,
	}
	c.messages <- message

	// put our listener on the queue and listen for it
	c.lock.Lock()
	results := make(chan *Response, 1)
	c.waiting[message.ID] = results
	c.lock.Unlock()

	select {
	case result := <-results:
		return result.Payload, result.Error

	case <-c.context.Done():
		return nil, ErrStopped
	}
}
