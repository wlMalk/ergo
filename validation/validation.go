package validation

import (
	"net/http"
)

type Validator interface {
	Validate(Valuer, Requester) error
}

type ValidatorFunc func(Valuer, Requester) error

func (f ValidatorFunc) Validate(v Valuer, r Requester) error {
	return f(v, r)
}

type Requester interface {
	Req() *http.Request
	Param(string) Valuer
	ParamOk(string) (Valuer, bool)
	Params(...string) map[string]Valuer
}

type Valuer interface {
	Name() string
	Value() interface{}
	As() int
	Int() int
	Int64() int64
	Float() float32
	Float64() float64
	String() string
	Bool() bool
	// IntE() (int, error)
	// Int64E() (int64, error)
	// FloatE() (float32, error)
	// Float64E() (float64, error)
	// BoolE() (bool, error)
}
