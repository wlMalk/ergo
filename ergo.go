package ergo

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

