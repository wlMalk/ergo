package ergo

type Context struct {
	Ergo      Ergoer
	Operation Operationer
	Response  *Response
	Request   *Request
}

func NewContext(res *Response, req *Request) *Context {
	return &Context{
		Response: res,
		Request:  req,
	}
}
