package ergo

// Route

type Route struct {
	parent      *Route
	name        string
	path        string
	description string
	routes      map[string]*Route
	params          map[string]*Param
	notFoundHandler Handler
}

func NewRoute(path string) *Route {
	return &Route{
		path:     path,
		routes:   map[string]*Route{},
		params:     map[string]*Param{},
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
	return r.path
}

// GetFullPath returns the absolute path of the route.
func (r *Route) GetFullPath() string {
	if r.parent != nil {
		return r.parent.GetFullPath() + r.path
	}
	return r.path
}

// GetRoutes returns the map of child routes.
func (r *Route) GetRoutes() map[string]*Route {
	return r.routes
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
