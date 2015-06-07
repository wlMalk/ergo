package ergo

import (
	"strings"

	"github.com/wlMalk/ergo/constants"
)

var (
	MaxMemory int64 = 32 << 20 // 32 MB
)


type OperationMap map[string]*Operation

func (om OperationMap) GetOperation(method string) (*Operation, bool) {
	o, ok := om[strings.ToUpper(method)]
	if !ok {
		o, ok = om[""]
	}
	return o, ok
}

type Operationer interface {
	GetSchemes() []string
	GetConsumes() []string
	GetProduces() []string
	GetName() string
	GetDescription() string
	GetMethod() string
	GetRoute() Router
}

// Operation

type Operation struct {
	root          *Ergo
	route         *Route
	method        string
	name          string
	description   string
	handler       Handler
	params        map[string]*Param
	schemes       []string
	consumes      []string
	produces      []string
	containsFiles bool
	bodyParams    bool
}

func NewOperation(handler Handler) *Operation {
	if handler == nil {
		panic("Handler cannot be nil")
	}
	return &Operation{
		handler: handler,
		params:  map[string]*Param{},
	}
}

func ANY(function HandlerFunc) *Operation {
	return HandleANY(HandlerFunc(function))
}

func HandleANY(handler Handler) *Operation {
	return NewOperation(handler).Method(constants.METHOD_ANY)
}

func GET(function HandlerFunc) *Operation {
	return HandleGET(HandlerFunc(function))
}

func HandleGET(handler Handler) *Operation {
	return NewOperation(handler).Method(constants.METHOD_GET)
}

func POST(function HandlerFunc) *Operation {
	return HandlePOST(HandlerFunc(function))
}

func HandlePOST(handler Handler) *Operation {
	return NewOperation(handler).Method(constants.METHOD_POST)
}

func PUT(function HandlerFunc) *Operation {
	return HandlePUT(HandlerFunc(function))
}

func HandlePUT(handler Handler) *Operation {
	return NewOperation(handler).Method(constants.METHOD_PUT)
}

func DELETE(function HandlerFunc) *Operation {
	return HandleDELETE(HandlerFunc(function))
}

func HandleDELETE(handler Handler) *Operation {
	return NewOperation(handler).Method(constants.METHOD_DELETE)
}

// Name sets the name of the operation.
func (o *Operation) Name(name string) *Operation {
	o.name = name
	return o
}

// GetName returns the name of the operation.
func (o *Operation) GetName() string {
	return o.name
}

// GetName returns the description of the operation.
func (o *Operation) GetDescription() string {
	return o.description
}

// Description sets the description of the operation.
func (o *Operation) Description(description string) *Operation {
	o.description = description
	return o
}

func (o *Operation) Method(method string) *Operation {
	if method == constants.METHOD_ANY ||
		method == constants.METHOD_GET ||
		method == constants.METHOD_POST ||
		method == constants.METHOD_PUT ||
		method == constants.METHOD_DELETE {
		o.method = method
	}
	return o
}

func (o *Operation) GetMethod() string {
	return o.method
}

// Schemes is not additive, meaning that it'll reset the schemes
// already defined with what it's been given if they are valid.
func (o *Operation) Schemes(s ...string) *Operation {
	schemes(o, s)
	return o
}

func (o *Operation) Consumes(mimes ...string) *Operation {
	consumes(o, mimes)
	return o
}

func (o *Operation) Produces(mimes ...string) *Operation {
	produces(o, mimes)
	return o
}

func (o *Operation) GetSchemes() []string {
	return o.schemes
}

func (o *Operation) GetConsumes() []string {
	return o.consumes
}

func (o *Operation) GetProduces() []string {
	return o.produces
}

func (o *Operation) GetRoute() Router {
	return o.route
}

func (o *Operation) Params(params ...*Param) *Operation {
	addParams(o, params...)
	return o
}

func (o *Operation) GetParams() map[string]*Param {
	return o.params
}

func (o *Operation) GetParamsSlice() []*Param {
	var params []*Param
	for _, p := range o.params {
		params = append(params, p)
	}
	return params
}

func (o *Operation) ResetParams(params ...*Param) *Operation {
	o.setParamsSlice(params...)
	return o
}

func (o *Operation) SetParams(params map[string]*Param) *Operation {
	o.setParams(params)
	return o
}

func (o *Operation) IgnoreParams(params ...string) *Operation {
	ignoreParams(o, params...)
	return o
}

func (o *Operation) IgnoreParamsBut(params ...string) *Operation {
	ignoreParamsBut(o, params...)
	return o
}

func (o *Operation) ServeHTTP(res *Response, req *Request) {
}

func (o *Operation) setSchemes(schemes []string) {
	o.schemes = schemes
}

func (o *Operation) setConsumes(mimes []string) {
	o.consumes = mimes
}

func (o *Operation) setProduces(mimes []string) {
	o.produces = mimes
}

func (o *Operation) setParams(params map[string]*Param) {
	if params == nil {
		params = make(map[string]*Param)
	}
	o.params = params
}

func (o *Operation) setParamsSlice(params ...*Param) {
	paramsMap := map[string]*Param{}
	for _, p := range params {
		o.params[p.name] = p
	}
	o.setParams(paramsMap)
}
