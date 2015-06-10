package validation

import (
	"strconv"

	"github.com/wlMalk/ergo/constants"
)

type Value struct {
	name             string
	value            interface{}
	strValue         string
	as               int
	strMultipleValue []string
	multiple         bool
	from             string
}

func NewValue(name string, value string, from string) *Value {
	return &Value{
		name:     name,
		strValue: value,
		from:     from,
		as:       constants.PARAM_STRING,
	}
}

func NewMultipleValue(name string, value []string, from string) *Value {
	return &Value{
		name:             name,
		strMultipleValue: value,
		multiple:         true,
		from:             from,
		as:               constants.PARAM_STRING,
	}
}

func (v *Value) Name() string {
	return v.name
}

func (v *Value) As() int {
	return v.as
}

func (v *Value) Value() interface{} {
	if v.value == nil {
		switch v.as {
		case constants.PARAM_STRING:
			v.String()
		case constants.PARAM_INT:
			v.Int()
		case constants.PARAM_INT64:
			v.Int64()
		case constants.PARAM_FLOAT:
			v.Float()
		case constants.PARAM_FLOAT64:
			v.Float64()
		case constants.PARAM_BOOL:
			v.Bool()
		}
	}
	return v.value
}

func (v *Value) RawString() string {
	return v.strValue
}

func (v *Value) String() string {
	if v.value != nil {
		rv := v.value.(string)
		return rv
	}
	if v.as == constants.PARAM_STRING {
		v.value = v.strValue
		return v.strValue
	}
	return ""
}

func (v *Value) Bool() bool {
	if v.value != nil {
		rv := v.value.(bool)
		return rv
	}
	if v.as == constants.PARAM_BOOL {
		rv, _ := strconv.ParseBool(v.strValue)
		v.value = rv
		return rv
	}
	return false
}

func (v *Value) Int() int {
	if v.value != nil {
		rv := v.value.(int)
		return rv
	}
	if v.as == constants.PARAM_INT {
		rv, _ := strconv.Atoi(v.strValue)
		v.value = rv
		return rv
	}
	return 0
}

func (v *Value) Int64() int64 {
	return 0
}

func (v *Value) Float() float32 {
	return 0
}

func (v *Value) Float64() float64 {
	return 0
}
