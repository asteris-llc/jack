package jack

type Request struct {
	ID     string        `json:"id"`
	Method string        `json:"method"`
	Args   []interface{} `json:"args"`
}

type Response struct {
	ID      string      `json:"id"`
	Payload interface{} `json:"payload"`
	Error   string      `json:"error"` // TODO: why can't encoding/json serialize an error?
}
