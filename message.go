package jack

type Message struct {
	ID      string
	Method  string
	Payload []interface{}
}
