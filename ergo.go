package ergo

type Ergo struct {
	*Route
}

func New(path string) *Ergo {
	return &Ergo{
		Route: NewRoute(path),
	}
}

