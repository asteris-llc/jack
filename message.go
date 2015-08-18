package jack

type Request struct {
	ID     string
	Method string
	Args   []interface{}
}

type Response struct {
	ID      string
	Payload interface{}
	Error   error
}
