package ergo

type Handler interface {
	ServeHTTP(*Response, *Request)
}

type HandlerFunc func(*Response, *Request)

func (f HandlerFunc) ServeHTTP(w *Response, r *Request) {
	f(w, r)
}

// Operation

type Operation struct {
	method      string
	name        string
	description string
	handler     Handler
}

func NewOperation(handler Handler) *Operation {
	return &Operation{
		handler: handler,
	}
}

