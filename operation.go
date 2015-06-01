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
	produces    []string
}

func NewOperation(handler Handler) *Operation {
	return &Operation{
		handler: handler,
	}
}

func ANY(function HandlerFunc) *Operation {
	return HandleANY(HandlerFunc(function))
}

func HandleANY(handler Handler) *Operation {
	return NewOperation(handler).Method(METHOD_ANY)
}

func GET(function HandlerFunc) *Operation {
	return HandleGET(HandlerFunc(function))
}

func HandleGET(handler Handler) *Operation {
	return NewOperation(handler).Method(METHOD_GET)
}

func POST(function HandlerFunc) *Operation {
	return HandlePOST(HandlerFunc(function))
}

func HandlePOST(handler Handler) *Operation {
	return NewOperation(handler).Method(METHOD_POST)
}

func PUT(function HandlerFunc) *Operation {
	return HandlePUT(HandlerFunc(function))
}

func HandlePUT(handler Handler) *Operation {
	return NewOperation(handler).Method(METHOD_PUT)
}

func DELETE(function HandlerFunc) *Operation {
	return HandleDELETE(HandlerFunc(function))
}

func HandleDELETE(handler Handler) *Operation {
	return NewOperation(handler).Method(METHOD_DELETE)
}

func (o *Operation) Description(description string) *Operation {
	o.description = description
	return o
}

func (o *Operation) GetDescription() string {
	return o.description
}

func (o *Operation) Method(method string) *Operation {
	o.method = method
	return o
}

func (o *Operation) GetMethod() string {
	return o.method
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

func (o *Operation) Produces(mimes ...string) *Operation {
	produces(o, mimes)
	return o
}

func (o *Operation) GetSchemes() []string {
	return o.schemes
}

func (o *Operation) GetConsumes() []string {
	return o.consumes
}

func (o *Operation) GetProduces() []string {
	return o.produces
}

func (o *Operation) setSchemes(schemes []string) {
	o.schemes = schemes
}

func (o *Operation) setConsumes(mimes []string) {
	o.consumes = mimes
}

func (o *Operation) setProduces(mimes []string) {
	o.produces = mimes
}
