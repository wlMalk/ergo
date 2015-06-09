package ergo

import (
	"net/http"

	"github.com/wlMalk/ergo/constants"
)

type Wrapper interface {
	Match(*http.Request) http.Handler
	Set([]*Operation)
}

var (
	DefSchemes  = []string{constants.SCHEME_HTTP}
	DefConsumes = []string{constants.MIME_JSON}
	DefProduces = []string{constants.MIME_JSON}
)

type Ergoer interface {
	GetSchemes() []string
	GetConsumes() []string
	GetProduces() []string
	NotFound(*Context)
	MethodNotAllowed(*Route, *Context)
	Err(error, *Response, *Request)
	Panic(*Context)
}

// Ergo

type Ergo struct {
	root   *Route
	router Wrapper

	schemes  []string
	consumes []string
	produces []string

	operations []*Operation

	NotFoundHandler      Handler
	MethodNotAllowedFunc MethodNotAllowedFunc
	ErrHandler           ErrHandler
	PanicHandler         Handler
}

func New() *Ergo {
	e := &Ergo{
		root:                 NewRoute(""),
		schemes:              DefSchemes,
		consumes:             DefConsumes,
		produces:             DefProduces,
		NotFoundHandler:      DefNotFoundHandler,
		MethodNotAllowedFunc: DefMethodNotAllowedFunc,
		ErrHandler:           DefErrHandler,
		PanicHandler:         DefPanicHandler,
	}
	return e
}

func NewWith(w Wrapper) *Ergo {
	e := New()
	e.Router(w)
	return e
}

func (e *Ergo) New(path string) *Route {
	return e.root.New(path)
}

// Use allows setting middleware for every operation in Ergo
func (e *Ergo) Use(middleware ...Middleware) *Ergo {
	e.root.Use(middleware...)
	return e
}

func (e *Ergo) UseFunc(middleware ...MiddlewareFunc) *Ergo {
	e.root.UseFunc(middleware...)
	return e
}

func (e *Ergo) Schemes(s ...string) *Ergo {
	schemes(e, s)
	return e
}

func (e *Ergo) Consumes(mimes ...string) *Ergo {
	consumes(e, mimes)
	return e
}

func (e *Ergo) Produces(mimes ...string) *Ergo {
	produces(e, mimes)
	return e
}

func (e *Ergo) Params(params ...*Param) *Ergo {
	addParams(e.root, params)
	return e
}

// Router uses a router that implement Router interface
// as the main router.
func (e *Ergo) Router(w Wrapper) {
	e.router = w
}

// GetSchemes returns the default schemes.
func (e *Ergo) GetSchemes() []string {
	return e.schemes
}

// GetConsumes returns the default consumable content types.
func (e *Ergo) GetConsumes() []string {
	return e.consumes
}

// GetProduces returns the default producible content types.
func (e *Ergo) GetProduces() []string {
	return e.produces
}

func (e *Ergo) setSchemes(schemes []string) {
	e.schemes = schemes
}

func (e *Ergo) setConsumes(consumes []string) {
	e.consumes = consumes
}

func (e *Ergo) setProduces(produces []string) {
	e.produces = produces
}

func (e *Ergo) Prepare() error {
	e.PrepareRouter()
	return nil
}

func (e *Ergo) PrepareRouter() {
	for _, r := range e.root.GetAllRoutes() {
		ops := r.GetOperations()
		var ge, po, pu, de bool
		for _, o := range ops {
			o.Handler = getHandler(o.Handler, append(o.middleware, MiddlewareFunc(o.Validate)))
			o.ergo = e
			switch o.method {
			case constants.METHOD_GET:
				ge = true
			case constants.METHOD_POST:
				po = true
			case constants.METHOD_PUT:
				pu = true
			case constants.METHOD_DELETE:
				de = true
			}
		}
		mnf := e.MethodNotAllowedFunc(r)
		if !ge {
			o := HandleGET(mnf)
			o.route = r
			ops = append(ops, o)
		}
		if !po {
			o := HandlePOST(mnf)
			o.route = r
			ops = append(ops, o)
		}
		if !pu {
			o := HandlePUT(mnf)
			o.route = r
			ops = append(ops, o)
		}
		if !de {
			o := HandleDELETE(mnf)
			o.route = r
			ops = append(ops, o)
		}
		e.operations = append(e.operations, ops...)
	}
	e.router.Set(e.operations)
}

func (e *Ergo) NotFound(ctx *Context) {
	e.NotFoundHandler.ServeHTTP(ctx)
}

func (e *Ergo) MethodNotAllowed(r *Route, ctx *Context) {
	e.MethodNotAllowedFunc(r).ServeHTTP(ctx)
}

func (e *Ergo) Err(err error, res *Response, req *Request) {
	e.ErrHandler.ServeHTTP(err, res, req)
}

func (e *Ergo) Panic(ctx *Context) {
	e.PanicHandler.ServeHTTP(ctx)
}

func (e *Ergo) Run(address string) error {
	err := e.Prepare()
	if err != nil {
		return err
	}
	return http.ListenAndServe(address, e)
}

func (e *Ergo) RunTLS(addr, certFile, keyFile string) error {
	err := e.Prepare()
	if err != nil {
		return err
	}
	return http.ListenAndServeTLS(addr, certFile, keyFile, e)
}

func (e *Ergo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := preparePath(r.URL.Path)
	if r.URL.Path != path && r.Method == "GET" {
		r.URL.Path = path
		http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
		return
	}

	handler := e.router.Match(r)
	if handler == nil {
		// not found
		return
	}

	handler.ServeHTTP(w, r)
}
