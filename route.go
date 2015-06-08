package ergo

import (
	"strings"
)

type Router interface {
}

// Route

type Route struct {
	ergo       *Ergo
	parent     *Route
	path       string
	routes     []*Route
	params     map[string]*Param
	operations []*Operation
}

func NewRoute(path string) *Route {
	return &Route{
		path:   strings.ToLower(preparePath(path)),
		params: map[string]*Param{},
	}
}

// GetParent returns the parent of the route.
func (r *Route) GetParent() *Route {
	return r.parent
}

// GetPath returns the relative path of the route to the
// parent route.
func (r *Route) GetPath() string {
	if r.parent != nil {
		return "/" + r.path
	}
	return ""
}

// GetFullPath returns the absolute path of the route.
func (r *Route) GetFullPath() string {
	if r.parent != nil {
		return r.parent.GetFullPath() + r.GetPath()
	}
	return r.GetPath()
}

// New creates a route with the provided path and adds it
// to r then returns a pointer to it.
func (r *Route) New(path string) *Route {
	route := NewRoute(path)
	r.addRoute(route)
	return route
}

// AddRoute copies and add the given route then returns
// a pointer to the added one.
func (r *Route) AddRoute(route *Route) *Route {
	nroute := route.Copy()
	r.addRoute(nroute)
	return nroute
}

// AddRoutes copies and add every Route in routes.
func (r *Route) AddRoutes(routes ...*Route) *Route {
	for _, nr := range routes {
		r.AddRoute(nr)
	}
	return r
}

// GetRoutes returns the slice of child routes.
func (r *Route) GetRoutes() []*Route {
	return r.routes
}

// GetRoutes returns the slice of all nested routes under r.
func (r *Route) GetAllRoutes() []*Route {
	var routes []*Route
	for _, route := range r.routes {
		routes = append(routes, route.GetAllRoutes()...)
	}
	return routes
}

// SetRoutes replaces the routes slice with the given one.
func (r *Route) SetRoutes(routes []*Route) *Route {
	r.routes = routes
	return r
}

// Params add the given params to the params map in the route.
// No two params can have the same name, even if the were
// in different places.
func (r *Route) Params(params ...*Param) *Route {
	addParams(r, params...)
	return r
}

func (r *Route) GetParams() map[string]*Param {
	return r.params
}

func (r *Route) GetParamsSlice() []*Param {
	var params []*Param
	for _, p := range r.params {
		params = append(params, p)
	}
	return params
}

func (r *Route) ResetParams(params ...*Param) *Route {
	r.setParamsSlice(params...)
	return r
}

func (r *Route) SetParams(params map[string]*Param) *Route {
	r.setParams(params)
	return r
}

func (r *Route) IgnoreParams(params ...string) *Route {
	ignoreParams(r, params...)
	return r
}

func (r *Route) IgnoreParamsBut(params ...string) *Route {
	ignoreParamsBut(r, params...)
	return r
}

func (r *Route) GET(function HandlerFunc) *Operation {
	return r.HandleGET(HandlerFunc(function))
}

func (r *Route) HandleGET(handler Handler) *Operation {
	operation := HandleGET(handler)
	setOperation(r, operation)
	r.operations[constants.METHOD_GET] = operation
	return operation
}

func (r *Route) POST(function HandlerFunc) *Operation {
	return r.HandlePOST(HandlerFunc(function))
}

func (r *Route) HandlePOST(handler Handler) *Operation {
	operation := HandlePOST(handler)
	setOperation(r, operation)
	r.operations[constants.METHOD_POST] = operation
	return operation
}

func (r *Route) PUT(function HandlerFunc) *Operation {
	return r.HandlePUT(HandlerFunc(function))
}

func (r *Route) HandlePUT(handler Handler) *Operation {
	operation := HandlePUT(handler)
	setOperation(r, operation)
	r.operations[constants.METHOD_PUT] = operation
	return operation
}

func (r *Route) DELETE(function HandlerFunc) *Operation {
	return r.HandleDELETE(HandlerFunc(function))
}

func (r *Route) HandleDELETE(handler Handler) *Operation {
	operation := HandleDELETE(handler)
	setOperation(r, operation)
	r.operations[constants.METHOD_DELETE] = operation
	return operation
}

// Operations does not alter the given operations in any way,
// it does not even add route parameters.
func (r *Route) Operations(operations ...*Operation) *Route {
	for _, o := range operations {
		r.operations[o.method] = o
	}
	return r
}

func (r *Route) GetOperations() OperationMap {
	return r.operations
}

func (r *Route) GetOperationers() OperationerMap {
	ops := OperationerMap{}
	for m, o := range r.operations {
		ops[m] = o
	}
	return ops
}

// returns Ergo object with only a set of methods exposed
func (r *Route) Ergo() Ergoer {
	return r.ergo
}

func (r *Route) Match(path string) (*Route, string) {
	if r.parent != nil {
		mat, rem, par := r.match(path)
		if !mat {
			return nil, ""
		}
		if rem == "" {
			return r, par
		}
		return r.subMatch(rem)
	}
	return r.subMatch(path)
}

func (r *Route) MatchURL(u *url.URL) (*Route, string) {
	return r.Match(u.Path[:len(u.Path)+1])
}

func (r *Route) ServeHTTP(res *Response, req *Request) {
	// validate the params with all the matching routes
	o, ok := r.operations.GetOperation(req.Method)
	if !ok {
		r.ergo.MethodNotAllowed(r, res, req)
		return
	}
	req.operation = o
	o.ServeHTTP(res, req)
}

// Copy returns a pointer to a copy of the route.
// It does not copy parent, operations, nor deep-copy the params.
func (r *Route) Copy() *Route {
	route := NewRoute(r.path)
	for _, cr := range r.routes {
		route.AddRoute(cr)
	}
	route.params = r.params
	return route
}

func (r *Route) addRoute(route *Route) {
	_, ok := r.routes[route.path]
	if ok {
		panic(fmt.Sprintf("A route with the path \"%s\" already exists.", route.path))
	}
	route.ergo = r.ergo
	route.parent = r
	setParamer(r, route)
	r.routesSlice = append(r.routesSlice, route.path)
	r.routes[route.path] = route
}

func (r *Route) setParams(params map[string]*Param) {
	if params == nil {
		params = make(map[string]*Param)
	}
	r.params = params
}

func (r *Route) setParamsSlice(params ...*Param) {
	paramsMap := map[string]*Param{}
	for _, p := range params {
		r.params[p.name] = p
	}
	r.setParams(paramsMap)
}

func (r *Route) match(path string) (bool, string, string) {
	return match(r.path, path)
}

func (r *Route) subMatch(path string) (*Route, string) {
	for _, routepath := range r.routesSlice {
		nr, par := r.routes[routepath].Match(path)
		if nr != nil {
			return nr, par
		}
	}
	return nil, ""
}


