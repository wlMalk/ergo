package ergo

import (
	"strconv"

	"github.com/wlMalk/ergo/validation"
)

// Param

type Param struct {
	name        string
	description string
	validators  []validation.Validator
	def         interface{}
	as          int
	required    bool
	multiple    bool
	strSep      string
	file        bool
	inPath      bool
	inQuery     bool
	inHeader    bool
	inBody      bool
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
	p.required = true
	return p
}

func (p *Param) File() *Param {
	p.file = true
	return p
}

func (p *Param) Multiple() *Param {
	p.multiple = true
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

// Validate returns the first error it encountered
func (p *Param) Validate(pv validation.Valuer, r *Request) error {
	for _, v := range p.validators {
		err := v.Validate(pv, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateAll returns all the errors it encountered
func (p *Param) ValidateAll(pv validation.Valuer, r *Request) []error {
	var errs []error
	for _, v := range p.validators {
		err := v.Validate(pv, r)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

type ParamValue struct {
	name             string
	value            interface{}
	strValue         string
	as               int
	strMultipleValue []string
	multiple         bool
	from             string
}

func NewParamValue(name string, value string, from string) *ParamValue {
	return &ParamValue{
		name:     name,
		strValue: value,
		from:     from,
	}
}

func NewMultipleParamValue(name string, value []string, from string) *ParamValue {
	return &ParamValue{
		name:             name,
		strMultipleValue: value,
		multiple:         true,
		from:             from,
	}
}

func (pv *ParamValue) Name() string {
	return pv.name
}

func (pv *ParamValue) As() int {
	return pv.as
}

func (pv *ParamValue) Value() interface{} {
	if pv.value == nil {
		return pv.strValue
	}
	return pv.value
}

func (pv *ParamValue) String() string {
	return pv.strValue
}

func (pv *ParamValue) Bool() bool {
	if pv.value != nil {
		v := pv.value.(bool)
		return v
	}
	if pv.as == PARAM_BOOL {
		v, _ := strconv.ParseBool(pv.strValue)
		pv.value = v
		return v
	}
	return false
}

func (pv *ParamValue) Int() int {
	if pv.value != nil {
		v := pv.value.(int)
		return v
	}
	if pv.as == PARAM_INT {
		v, _ := strconv.Atoi(pv.strValue)
		pv.value = v
		return v
	}
	return 0
}

func (pv *ParamValue) Int64() int64 {
	return 0
}

func (pv *ParamValue) Float() float32 {
	return 0
}

func (pv *ParamValue) Float64() float64 {
	return 0
}
