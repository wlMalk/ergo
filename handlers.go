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

type ErrHandler interface {
	ServeHTTP(error, *Response, *Request)
}
