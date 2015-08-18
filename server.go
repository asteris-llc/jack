package jack

import (
	"errors"
	"reflect"
)

var (
	ErrNoSuchMethod = errors.New("no such method")
)

type Server struct {
	value reflect.Value
}

func NewServer(instance interface{}) *Server {
	return &Server{
		value: reflect.ValueOf(instance),
	}
}

func (s *Server) Run() error {
	return nil
}

func (s *Server) handleCall(message *Request) *Response {
	method := s.value.MethodByName(message.Method)
	if !method.IsValid() {
		return &Response{
			ID:    message.ID,
			Error: ErrNoSuchMethod,
		}
	}

	values := []reflect.Value{}
	for _, arg := range message.Args {
		values = append(values, reflect.ValueOf(arg))
	}

	result := method.Call(values)
	var err error = nil
	if len(result) > 1 {
		maybeErr, ok := result[len(result)].Interface().(error)
		if ok {
			err = maybeErr
			result = result[0 : len(result)-1]
		}
	}

	return &Response{
		ID:      message.ID,
		Payload: result,
		Error:   err,
	}
}
