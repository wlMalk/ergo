package ergo

import (
	"net/http"
)

type Response struct {
	http.ResponseWriter
	StatusCode    int
	contentLength int
	indent        bool
}

func NewResponse(httpResponse http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: httpResponse,
		StatusCode:     http.StatusOK,
	}
}
