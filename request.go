package ergo

import (
	"net/http"

	"github.com/wlMalk/ergo/validation"
)

type Request struct {
	*http.Request
	Input      map[string]validation.Valuer
	pathParams map[string]string
	operation  *Operation // Operation object
}

func NewRequest(httpRequest *http.Request) *Request {
	return &Request{
		Request:    httpRequest,
		pathParams: map[string]string{},
		Input:      map[string]validation.Valuer{},
	}
}

// Req returns the request.
func (req *Request) Req() *http.Request {
	return req.Request
}

// Param returns the input parameter value by its name.
func (req *Request) Param(name string) validation.Valuer {
	return req.Input[name]
}

// ParamOk returns the input parameter value by its name.
func (req *Request) ParamOk(name string) (validation.Valuer, bool) {
	p, ok := req.Input[name]
	return p, ok
}

// Params returns a map of input parameters values by their names.
// If no names given then it returns r.Input
func (req *Request) Params(names ...string) map[string]validation.Valuer {
	if len(names) == 0 {
		return req.Input
	}
	params := map[string]validation.Valuer{}
	for _, n := range names {
		p, ok := req.Input[n]
		if !ok {
			continue
		}
		params[n] = p
	}
	return params
}

func (req *Request) GetOperation() Operationer {
	return req.operation
}

func (req *Request) SetPathParams(params map[string]string) {
	req.pathParams = params
}
