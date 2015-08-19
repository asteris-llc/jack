package jack

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"reflect"
)

var (
	ErrNoSuchMethod = errors.New("no such method")
	ErrBadArity     = errors.New("bad arity")
	ErrBadType      = errors.New("bad type")
)

type Server struct {
	value reflect.Value
}

func NewServer(instance interface{}) *Server {
	return &Server{
		value: reflect.ValueOf(instance),
	}
}

func (s *Server) Start() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		req := &Request{}
		err := json.Unmarshal(scanner.Bytes(), req)
		if err != nil {
			return err
		}

		resp, err := json.Marshal(s.handleCall(req))
		if err != nil {
			return err
		}
		os.Stdout.Write(resp)
		os.Stdout.WriteString("\n")
	}

	return scanner.Err()
}

func (s *Server) handleCall(message *Request) *Response {
	method := s.value.MethodByName(message.Method)
	if !method.IsValid() {
		return &Response{
			ID:    message.ID,
			Error: ErrNoSuchMethod.Error(),
		}
	}

	values := []reflect.Value{}
	for _, arg := range message.Args {
		values = append(values, reflect.ValueOf(arg))
	}

	meta := method.Type()
	args := []reflect.Value{}

	if len(values) != meta.NumIn() {
		return &Response{
			ID:    message.ID,
			Error: ErrBadArity.Error(),
		}
	}

	for i := 0; i < meta.NumIn(); i++ {
		should := meta.In(i)
		value := values[i]

		if value.Type().ConvertibleTo(should) {
			args = append(args, value.Convert(should))
		} else {
			return &Response{
				ID:    message.ID,
				Error: ErrBadType.Error(),
			}
		}
	}

	result := method.Call(args)

	serializable := []interface{}{}
	for _, value := range result {
		serializable = append(serializable, value.Interface())
	}

	return &Response{
		ID:      message.ID,
		Payload: serializable,
		Error:   "",
	}
}
