package validation

import (
	"github.com/wlMalk/ergo/constants"
)

// Param

type Param struct {
	name        string
	description string
	validators  []Validator
	def         interface{}
	as          int
	IsRequired  bool
	IsMultiple  bool
	strSep      string
	IsFile      bool
	IsInPath    bool
	IsInQuery   bool
	IsInHeader  bool
	IsInBody    bool
}

func NewParam(name string) *Param {
	return &Param{
		name: name,
	}
}

func PathParam(name string) *Param {
	return NewParam(name).In(constants.IN_PATH)
}

func QueryParam(name string) *Param {
	return NewParam(name).In(constants.IN_QUERY)
}

func HeaderParam(name string) *Param {
	return NewParam(name).In(constants.IN_HEADER)
}

func BodyParam(name string) *Param {
	return NewParam(name).In(constants.IN_BODY)
}

func (p *Param) Name(name string) *Param {
	p.name = name
	return p
}

func (p *Param) GetName() string {
	return p.name
}

func (p *Param) Description(description string) *Param {
	p.description = description
	return p
}

func (p *Param) GetDescription() string {
	return p.description
}

func (p *Param) Required() *Param {
	p.IsRequired = true
	return p
}

func (p *Param) File() *Param {
	p.IsFile = true
	return p
}

func (p *Param) Multiple() *Param {
	p.IsMultiple = true
	return p
}

func (p *Param) As(as int) *Param {
	p.as = as
	return p
}

// If a Param is in path then it is required.
func (p *Param) In(in ...int) *Param {
	for _, i := range in {
		switch i {
		case constants.IN_PATH:
			p.IsInPath = true
			p.IsRequired = true
		case constants.IN_QUERY:
			p.IsInQuery = true
		case constants.IN_HEADER:
			p.IsInHeader = true
		case constants.IN_BODY:
			p.IsInBody = true
		}
	}
	return p
}

// Must sets the validators to use.
func (p *Param) Must(validators ...Validator) *Param {
	p.validators = validators
	return p
}

// Validate returns the first error it encountered
func (p *Param) Validate(v *Value, req Requester) error {
	for _, va := range p.validators {
		err := va.Validate(v, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateAll returns all the errors it encountered
func (p *Param) ValidateAll(v *Value, req Requester) []error {
	var errs []error
	for _, va := range p.validators {
		err := va.Validate(v, req)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (p *Param) Copy() *Param {
	param := NewParam(p.name)
	param.description = p.description
	param.def = p.def
	param.as = p.as
	param.IsRequired = p.IsRequired
	param.IsMultiple = p.IsMultiple
	param.strSep = p.strSep
	param.IsFile = p.IsFile
	param.IsInPath = p.IsInPath
	param.IsInQuery = p.IsInQuery
	param.IsInHeader = p.IsInHeader
	param.IsInBody = p.IsInBody
	param.validators = p.validators
	return param
}
