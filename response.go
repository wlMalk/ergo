package ergo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/wlMalk/ergo/constants"
)

type Response struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
	Encoding      string
	indent        bool
}

func NewResponse(httpResponse http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: httpResponse,
	}
}

func (res *Response) WriteEntity(status int, value interface{}) error {
	if value == nil {
		return nil
	}

	switch res.Encoding {
	case constants.MIME_JSON:
		return res.WriteAsJson(status, value)
	case constants.MIME_XML:
		return res.WriteAsXml(status, value, true)
	}

	return nil
}

func (res *Response) Stream(status int, d time.Duration, f func(int64) (interface{}, bool)) error {
	if f == nil {
		return nil
	}

	switch res.Encoding {
	case constants.MIME_JSON:
		return res.StreamAsJson(status, d, f)
	case constants.MIME_XML:
		return res.StreamAsXml(status, d, f)
	}

	return nil
}

func (res *Response) WriteAsXml(status int, value interface{}, writeHeader bool) error {
	var output []byte
	var err error

	if value == nil {
		return nil
	}
	if res.indent {
		output, err = xml.MarshalIndent(value, " ", " ")
	} else {
		output, err = xml.Marshal(value)
	}

	if err != nil {
		return res.WriteError(http.StatusInternalServerError, err)
	}
	res.Header().Set(constants.HEADER_ContentType, constants.MIME_XML)
	res.WriteHeader(status)
	if writeHeader {
		_, err = res.Write([]byte(xml.Header))
		if err != nil {
			return err
		}
	}
	if _, err = res.Write(output); err != nil {
		return err
	}
	return nil
}

func (res *Response) StreamAsXml(status int, d time.Duration, f func(int64) (interface{}, bool)) error {
	var i int64 = 0
	var err error
	for {
		e, stop := f(i)
		if e != nil {
			eerr := res.WriteAsXml(status, e, i == 0)
			if eerr != nil {
				err = eerr
				break
			}
			_, eerr = res.Write([]byte("\n"))
			res.Flush()
			i++
		}
		if stop {
			break
		}
		time.Sleep(d)
	}

	return err
}

func (res *Response) WriteAsJson(status int, value interface{}) error {
	return res.WriteJson(status, value, constants.MIME_JSON)
}

func (res *Response) StreamAsJson(status int, d time.Duration, f func(int64) (interface{}, bool)) error {
	var i int64 = 0
	var err error
	for {
		e, stop := f(i)
		if e != nil {
			eerr := res.WriteJson(status, e, constants.MIME_JSON)
			if eerr != nil {
				err = eerr
				break
			}
			_, eerr = res.Write([]byte("\n"))
			res.Flush()
			i++
		}
		if stop {
			break
		}
		time.Sleep(d)
	}

	return err
}

func (res *Response) WriteJson(status int, value interface{}, contentType string) error {
	var output []byte
	var err error

	if value == nil {
		return nil
	}
	if res.indent {
		output, err = json.MarshalIndent(value, " ", " ")
	} else {
		output, err = json.Marshal(value)
	}

	if err != nil {
		return res.WriteString(http.StatusInternalServerError, err.Error())
	}
	res.Header().Set(constants.HEADER_ContentType, contentType)
	res.WriteHeader(status)
	if _, err = res.Write(output); err != nil {
		return err
	}
	return nil
}

func (res *Response) WriteError(httpStatus int, err error) error {
	return res.WriteString(httpStatus, err.Error())
}

func (res *Response) WriteString(status int, str string) error {
	res.WriteHeader(status)
	if _, err := res.Write([]byte(str)); err != nil {
		return err
	}
	return nil
}

func (res *Response) WriteHeader(httpStatus int) {
	if res.statusCode == 0 {
		if httpStatus == 0 {
			httpStatus = http.StatusOK
		}
		res.statusCode = httpStatus
		res.ResponseWriter.WriteHeader(httpStatus)
	}
}

func (res *Response) StatusCode() int {
	if res.statusCode == 0 {
		return http.StatusOK
	}
	return res.statusCode
}

func (res *Response) Write(bytes []byte) (int, error) {
	written, err := res.ResponseWriter.Write(bytes)
	res.contentLength += written
	return written, err
}

func (res *Response) Indented(indent bool) {
	res.indent = indent
}

func (res *Response) ContentLength() int {
	return res.contentLength
}

func (res *Response) Flush() {
	if f, ok := res.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (res *Response) CloseNotify() <-chan bool {
	return res.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
