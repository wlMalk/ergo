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

// GetSchemes returns the default schemes passed from
// the parent.
func (r *Route) GetSchemes() []string {
	return r.schemes
}

func (r *Route) setSchemes(schemes []string) {
	r.schemes = schemes
}

