package ergo

// Param

type Param struct {
	name        string
	description string
	def      interface{}
	as       int
	required bool
	multiple bool
	strSep   string
	file     bool
	inPath   bool
	inQuery  bool
	inHeader bool
	inBody   bool
}

func NewParam(name string) *Param {
	return &Param{
		name: name,
	}
}

func PathParam(name string) *Param {
	return NewParam(name).In(IN_PATH)
}

func QueryParam(name string) *Param {
	return NewParam(name).In(IN_QUERY)
}

func HeaderParam(name string) *Param {
	return NewParam(name).In(IN_HEADER)
}

func BodyParam(name string) *Param {
	return NewParam(name).In(IN_BODY)
}

// If a Param is in path then it is required.
func (p *Param) In(in ...int) *Param {
	for _, i := range in {
		switch i {
		case IN_PATH:
			p.inPath = true
			p.required = true
		case IN_QUERY:
			p.inQuery = true
		case IN_HEADER:
			p.inHeader = true
		case IN_BODY:
			p.inBody = true
		}
	}
	return p
}

