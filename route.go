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

