package ergo

import (
	"net/http"
	"net/url"
)

type Ergo struct {
	*Route
}

func New(path string) *Ergo {
	return &Ergo{
		Route: NewRoute(path),
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

func (e *Ergo) FindRouteFromURL(urlObj *url.URL) (*Route, map[string]string) {

	return nil, nil
}

func (e *Ergo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	prepared := preparePath(r.URL.Path)
	if r.URL.Path != prepared && (r.Method == "GET" || r.Method == "") {
		r.URL.Path = prepared
		http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
		return
	}
	r.URL.Path = prepared
	route, pathParams := e.FindRouteFromURL(r.URL)
	if route != nil {
		req := NewRequest(r)
		req.route = route
		req.pathParams = pathParams
		res := NewResponse(w)
		route.ServeHTTP(res, req)
	}
}
