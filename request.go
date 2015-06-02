package ergo

import (
	"net/http"
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
