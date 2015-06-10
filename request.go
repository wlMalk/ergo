package ergo

import (
	"net/http"

	"github.com/wlMalk/ergo/validation"
)

type Request struct {
	*http.Request
	Input      map[string]*validation.Value
	PathParams map[string]string
}

func NewRequest(httpRequest *http.Request) *Request {
	return &Request{
		Request:    httpRequest,
		PathParams: map[string]string{},
		Input:      map[string]*validation.Value{},
	}
}

// Req returns the request.
func (req *Request) Req() *http.Request {
	return req.Request
}

// Param returns the input parameter value by its name.
func (req *Request) Param(name string) *validation.Value {
	return req.Input[name]
}

// ParamOk returns the input parameter value by its name.
func (req *Request) ParamOk(name string) (*validation.Value, bool) {
	p, ok := req.Input[name]
	return p, ok
}

// Params returns a map of input parameters values by their names.
// If no names given then it returns r.Input
func (req *Request) Params(names ...string) map[string]*validation.Value {
	if len(names) == 0 {
		return req.Input
	}
	params := map[string]*validation.Value{}
	for _, n := range names {
		p, ok := req.Input[n]
		if !ok {
			continue
		}
		params[n] = p
	}
	return params
}
