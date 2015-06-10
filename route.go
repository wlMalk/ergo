package ergo

import (
	"strings"

	"github.com/wlMalk/ergo/validation"
)

type Router interface {
	paramer
	GetPath() string
	GetFullPath() string
}

// Route

type Route struct {
	parent     *Route
	path       string
	routes     []*Route
	params     map[string]*validation.Param
	operations []*Operation
	middleware []Middleware
}

func NewRoute(path string) *Route {
	return &Route{
		path:   strings.ToLower(preparePath(path)),
		params: map[string]*validation.Param{},
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
		return r.path
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

// Use allows adding middleware for any operation in r
// Use is additive
func (r *Route) Use(middleware ...Middleware) *Route {
	r.middleware = append(r.middleware, middleware...)
	return r
}

func (r *Route) UseFunc(middleware ...MiddlewareFunc) *Route {
	for _, m := range middleware {
		r.Use(Middleware(m))
	}
	return r
}

func (r *Route) GetMiddleware() []Middleware {
	return r.middleware
}

// SetMiddleware will reset middleware with the given middleware
func (r *Route) SetMiddleware(middleware ...Middleware) *Route {
	r.middleware = middleware
	return r
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
		routes = append(routes, route)
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
func (r *Route) Params(params ...*validation.Param) *Route {
	addParams(r, params)
	return r
}

func (r *Route) GetParams() map[string]*validation.Param {
	return r.params
}

func (r *Route) GetParamsSlice() []*validation.Param {
	var params []*validation.Param
	for _, p := range r.params {
		params = append(params, p)
	}
	return params
}

func (r *Route) SetParamsSlice(params ...*validation.Param) *Route {
	r.setParamsSlice(params...)
	return r
}

func (r *Route) SetParams(params map[string]*validation.Param) *Route {
	r.setParams(params)
	return r
}

func (r *Route) IgnoreParams(params ...string) *Route {
	ignoreParams(r, params)
	return r
}

func (r *Route) IgnoreParamsBut(params ...string) *Route {
	ignoreParamsBut(r, params)
	return r
}

func (r *Route) GET(function HandlerFunc) *Operation {
	return r.HandleGET(HandlerFunc(function))
}

func (r *Route) HandleGET(handler Handler) *Operation {
	operation := HandleGET(handler)
	r.addOperation(operation)
	return operation
}

func (r *Route) POST(function HandlerFunc) *Operation {
	return r.HandlePOST(HandlerFunc(function))
}

func (r *Route) HandlePOST(handler Handler) *Operation {
	operation := HandlePOST(handler)
	r.addOperation(operation)
	return operation
}

func (r *Route) PUT(function HandlerFunc) *Operation {
	return r.HandlePUT(HandlerFunc(function))
}

func (r *Route) HandlePUT(handler Handler) *Operation {
	operation := HandlePUT(handler)
	r.addOperation(operation)
	return operation
}

func (r *Route) DELETE(function HandlerFunc) *Operation {
	return r.HandleDELETE(HandlerFunc(function))
}

func (r *Route) HandleDELETE(handler Handler) *Operation {
	operation := HandleDELETE(handler)
	r.addOperation(operation)
	return operation
}

// Operations does not alter the given operations in any way,
// it does not even add route parameters.
func (r *Route) Operations(operations ...*Operation) *Route {
	r.operations = operations
	return r
}

func (r *Route) GetAllOperations() []*Operation {
	var ops []*Operation
	for _, route := range r.routes {
		ops = append(ops, route.operations...)
		ops = append(ops, route.GetAllOperations()...)
	}
	return ops
}

func (r *Route) GetOperations() []*Operation {
	return r.operations
}

// Copy returns a pointer to a copy of the route.
// It does not copy parent, operations, nor deep-copy the params.
func (r *Route) Copy() *Route {
	route := NewRoute(r.path)
	for _, cr := range r.routes {
		route.AddRoute(cr)
	}
	route.params = r.params
	route.middleware = r.middleware
	return route
}

func (r *Route) ServeHTTP(ctx *Context) {
	for _, o := range r.operations {
		if o.method == ctx.Request.Method {
			o.ServeHTTP(ctx)
			return
		}
	}
	// method not allowed
	return
}

func (r *Route) addOperation(o *Operation) {
	o.route = r
	o.Params(r.GetParamsSlice()...)
	r.operations = append(r.operations, o)
	o.Use(r.middleware...)
}

func (r *Route) addRoute(route *Route) {
	route.parent = r
	route.Params(r.GetParamsSlice()...)
	r.routes = append(r.routes, route)
	route.Use(r.middleware...)
}

func (r *Route) setParams(params map[string]*validation.Param) {
	if params == nil {
		params = make(map[string]*validation.Param)
	}
	r.params = params
}

func (r *Route) setParamsSlice(params ...*validation.Param) {
	paramsMap := map[string]*validation.Param{}
	for _, p := range params {
		paramsMap[p.GetName()] = p
	}
	r.setParams(paramsMap)
}

