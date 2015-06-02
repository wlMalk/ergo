package ergo

import (
	"net/http"
)

type Response struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
	indent        bool
}

func NewResponse(httpResponse http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: httpResponse,
		statusCode:     http.StatusOK,
	}
}
