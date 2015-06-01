package ergo

// Route

type Route struct {
	parent      *Route
	name        string
	path        string
	description string
	routes      map[string]*Route
}

func NewRoute(path string) *Route {
	return &Route{
		path:     path,
		routes:   map[string]*Route{},
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

func (r *Route) setSchemes(schemes []string) {
	r.schemes = schemes
}

