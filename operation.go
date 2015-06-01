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
	schemes     []string
	consumes    []string
}

func NewOperation(handler Handler) *Operation {
	return &Operation{
		handler: handler,
	}
}

// Schemes is not additive, meaning that it'll reset the schemes
// already defined with what it's been given if they are valid.
func (o *Operation) Schemes(s ...string) *Operation {
	schemes(o, s)
	return o
}

func (o *Operation) Consumes(mimes ...string) *Operation {
	consumes(o, mimes)
	return o
}

func (o *Operation) GetSchemes() []string {
	return o.schemes
}

func (o *Operation) GetConsumes() []string {
	return o.consumes
}

func (o *Operation) setSchemes(schemes []string) {
	o.schemes = schemes
}

func (o *Operation) setConsumes(mimes []string) {
	o.consumes = mimes
}

