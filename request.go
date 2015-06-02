package ergo

import (
	"net/http"

	"github.com/wlMalk/ergo/validation"
)

type Request struct {
	*http.Request
	Input      map[string]validation.Valuer
	pathParams map[string]string
	route      *Route // route object that matched request
}

func NewRequest(httpRequest *http.Request) *Request {
	return &Request{
		Request: httpRequest,
		Input:   map[string]validation.Valuer{},
	}
}

// Req returns the request.
func (r *Request) Req() *http.Request {
	return r.Request
}

// Param returns the input parameter value by its name.
func (r *Request) Param(name string) validation.Valuer {
	return r.Input[name]
}

// ParamOk returns the input parameter value by its name.
func (r *Request) ParamOk(name string) (validation.Valuer, bool) {
	p, ok := r.Input[name]
	return p, ok
}

// Params returns a map of input parameters values by their names.
// If no names given then it returns r.Input
func (r *Request) Params(names ...string) map[string]validation.Valuer {
	if len(names) == 0 {
		return r.Input
	}
	params := map[string]validation.Valuer{}
	for _, n := range names {
		p, ok := r.Input[n]
		if !ok {
			continue
		}
		params[n] = p
	}
	return params
}
