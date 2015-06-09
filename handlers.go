package ergo

var (
	DefNotFoundHandler      Handler
	DefMethodNotAllowedFunc MethodNotAllowedFunc = MethodNotAllowedFunc(methodNotAllowedFunc)
	DefErrHandler           ErrHandler
	DefPanicHandler         Handler
)

type Handler interface {
	ServeHTTP(*Context)
}

type HandlerFunc func(*Context)

func (f HandlerFunc) ServeHTTP(ctx *Context) {
	f(ctx)
}

type NoCtxHandlerFunc func(*Response, *Request)

func (f NoCtxHandlerFunc) ServeHTTP(ctx *Context) {
	f(ctx.Response, ctx.Request)
}

type Middleware interface {
	Run(Handler) Handler
}

type MiddlewareFunc func(Handler) Handler

func (f MiddlewareFunc) Run(h Handler) Handler {
	return f(h)
}

type MethodNotAllowedFunc func(Router) Handler

func methodNotAllowedFunc(r Router) Handler {
	return HandlerFunc(func(ctx *Context) {
	})
}

type ErrHandler interface {
	ServeHTTP(error, *Response, *Request)
}

type ErrHandlerFunc func(error, *Response, *Request)

func (f ErrHandlerFunc) ServeHTTP(err error, res *Response, req *Request) {
	f(err, res, req)
}
