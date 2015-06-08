package ergo

var (
	defaultNotFoundHandler         Handler
	defaultMethodNotAllowedHandler MethodNotAllowedHandler
	defaultErrHandler              ErrHandler
	defaultPanicHandler            Handler
)

type Handler interface {
	ServeHTTP(*Response, *Request)
}

type HandlerFunc func(*Response, *Request)

func (f HandlerFunc) ServeHTTP(res *Response, req *Request) {
	f(res, req)
}

type MethodNotAllowedHandler interface {
	ServeHTTP(*Route, *Response, *Request)
}

type MethodNotAllowedHandlerFunc func(*Route, *Response, *Request)

func (f MethodNotAllowedHandlerFunc) ServeHTTP(r *Route, res *Response, req *Request) {
	f(r, res, req)
}

type ErrHandler interface {
	ServeHTTP(error, *Response, *Request)
}

type ErrHandlerFunc func(error, *Response, *Request)

func (f ErrHandlerFunc) ServeHTTP(err error, res *Response, req *Request) {
	f(err, res, req)
}
