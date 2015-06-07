package ergo

import (
	"net/http"
	"strings"
)


// Ergo

type Ergo struct {
	root   *Route
	router ExternalRouter

	schemes  []string
	consumes []string
	produces []string

	NotFoundHandler         Handler
	MethodNotAllowedHandler MethodNotAllowedHandler
	ErrHandler              ErrHandler
	PanicHandler            Handler
}

func New() *Ergo {
	e := &Ergo{
		NotFoundHandler:         defaultNotFoundHandler,
		MethodNotAllowedHandler: defaultMethodNotAllowedHandler,
		ErrHandler:              defaultErrHandler,
		PanicHandler:            defaultPanicHandler,
	}
	r := NewRoute("")
	r.ergo = e
	e.root = r
	return e
}
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
	addParams(e, params...)
	return e
}

func (e *Ergo) ResetParams(params ...*Param) *Ergo {
	e.setParamsSlice(params...)
	return e
}

func (e *Ergo) SetParams(params map[string]*Param) *Ergo {
	e.setParams(params)
	return e
}

func (e *Ergo) IgnoreParams(params ...string) *Ergo {
	ignoreParams(e, params...)
	return e
}

func (e *Ergo) IgnoreParamsBut(params ...string) *Ergo {
	ignoreParamsBut(e, params...)
	return e
}

func (e *Ergo) NotFoundHandler(h Handler) *Ergo {
	e.Route.notFoundHandler = h
	return e
}

func (e *Ergo) Prepare() error {
	return nil
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
	if r.URL.Path != "/"+path && r.Method == "GET" {
		r.URL.Path = "/" + path
		http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
		return
	}
	route, rp := e.Match(path)
	if route == nil {
		// not found
		return
	}

	req := NewRequest(r)
	req.route = route
	if len(rp) > 0 {
		ps := strings.Split(rp[:len(rp)], ";")
		for _, p := range ps {
			ci := strings.Index(p, ":")

			req.pathParams[p[:ci]] = p[ci+1:]
		}
	}
	res := NewResponse(w)
	route.ServeHTTP(res, req)
}
