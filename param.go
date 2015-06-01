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

