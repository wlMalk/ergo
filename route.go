package ergo

import (
	"fmt"
	"strings"
)

// Route

type Route struct {
	parent          *Route
	name            string
	path            string
	description     string
	routes          map[string]*Route
	indexSlice      []string
	indexMap        map[string]int
	params          map[string]*Param
	schemes         []string
	consumes        []string
	produces        []string
	operations      OperationMap
	notFoundHandler Handler
}

func NewRoute(path string) *Route {
	return &Route{
		path:       strings.ToLower(preparePath(path)),
		routes:     map[string]*Route{},
		params:     map[string]*Param{},
		indexMap:   map[string]int{},
		operations: OperationMap{},
	}
}

// GetParent returns the parent of the route.
func (r *Route) GetParent() *Route {
	return r.parent
}

// Name sets the name of the route.
// The route name is only used for documentation purposes.
func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}

// GetName returns the name of the route.
func (r *Route) GetName() string {
	return r.name
}

// GetName returns the description of the route.
func (r *Route) GetDescription() string {
	return r.description
}

// Description sets the description of the route.
func (r *Route) Description(description string) *Route {
	r.description = description
	return r
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

// GetRoutes returns the map of child routes.
func (r *Route) GetRoutes() map[string]*Route {
	return r.routes
}

// GetRoutesSlice returns a slice of child routes
// based on the order in which it was added.
func (r *Route) GetRoutesSlice() []*Route {
	var routes []*Route
	for _, s := range r.indexSlice {
		routes = append(routes, r.routes[s])
	}
	return routes
}

// SetRoutes replaces the routes map with the given one.
func (r *Route) SetRoutes(routes map[string]*Route) *Route {
	r.routes = routes
	return r
}

// GetSchemes returns the default schemes passed from
// the parent.
func (r *Route) GetSchemes() []string {
	return r.schemes
}

// GetConsumes returns the consumable content types
// passed from the parent.
func (r *Route) GetConsumes() []string {
	return r.consumes
}

// GetProduces returns the producible content types
// passed from the parent.
func (r *Route) GetProduces() []string {
	return r.produces
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

func (r *Route) ANY(function HandlerFunc) *Operation {
	return r.HandleANY(HandlerFunc(function))
}

func (r *Route) HandleANY(handler Handler) *Operation {
	operation := HandleANY(handler)
	setChild(r, operation)
	r.operations[METHOD_ANY] = operation
	return operation
}

func (r *Route) GET(function HandlerFunc) *Operation {
	return r.HandleGET(HandlerFunc(function))
}

func (r *Route) HandleGET(handler Handler) *Operation {
	operation := HandleGET(handler)
	setChild(r, operation)
	r.operations[METHOD_GET] = operation
	return operation
}

func (r *Route) POST(function HandlerFunc) *Operation {
	return r.HandlePOST(HandlerFunc(function))
}

func (r *Route) HandlePOST(handler Handler) *Operation {
	operation := HandlePOST(handler)
	setChild(r, operation)
	r.operations[METHOD_POST] = operation
	return operation
}

func (r *Route) PUT(function HandlerFunc) *Operation {
	return r.HandlePUT(HandlerFunc(function))
}

func (r *Route) HandlePUT(handler Handler) *Operation {
	operation := HandlePUT(handler)
	setChild(r, operation)
	r.operations[METHOD_PUT] = operation
	return operation
}

func (r *Route) DELETE(function HandlerFunc) *Operation {
	return r.HandleDELETE(HandlerFunc(function))
}

func (r *Route) HandleDELETE(handler Handler) *Operation {
	operation := HandleDELETE(handler)
	setChild(r, operation)
	r.operations[METHOD_DELETE] = operation
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

// NotFoundHandler sets the handler used when an operation is not found
// in the route, when a subroute could not be found.
func (r *Route) NotFoundHandler(h Handler) *Route {
	r.notFoundHandler = h
	return r
}

// GetNotFoundHandler returns the handler set in the route.
// If it is nil and t is true then it will try and look for handler in a parent.
func (r *Route) GetNotFoundHandler(t bool) Handler {
	if r.notFoundHandler == nil {
		if t {
			return r.parent.GetNotFoundHandler(t)
		}
		return nil
	}
	return r.notFoundHandler
}

func (r *Route) ServeHTTP(res *Response, req *Request) {
}

// Copy returns a pointer to a copy of the route.
// It does not copy parent, operations, nor deep-copy the params.
func (r *Route) Copy() *Route {
	route := NewRoute(r.path)
	route.name = r.name
	route.description = r.description
	for _, cr := range r.routes {
		route.AddRoute(cr)
	}
	route.params = r.params
	route.schemes = r.schemes
	route.consumes = r.consumes
	route.produces = r.produces
	route.notFoundHandler = r.notFoundHandler
	return route
}

func (r *Route) addRoute(route *Route) {
	_, ok := r.routes[route.path]
	if ok {
		panic(fmt.Sprintf("A route with the path \"%s\" already exists.", route.path))
	}
	route.parent = r
	setChild(r, route)
	r.indexSlice = append(r.indexSlice, route.path)
	r.indexMap[route.path] = len(r.indexSlice) - 1
	r.routes[route.path] = route
}

func (r *Route) setSchemes(schemes []string) {
	r.schemes = schemes
}

func (r *Route) setConsumes(consumes []string) {
	r.consumes = consumes
}

func (r *Route) setProduces(produces []string) {
	r.produces = produces
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
