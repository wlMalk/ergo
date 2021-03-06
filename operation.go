package ergo

import (
	"net/http"
	"strings"

	"github.com/wlMalk/ergo/constants"
	"github.com/wlMalk/ergo/validation"
)

var (
	MaxMemory int64 = 32 << 20 // 32 MB
)

type Operationer interface {
	paramer
	schemer
	consumer
	producer
	GetName() string
	GetDescription() string
	GetMethod() string
	GetRoute() Router
}

// Operation

type Operation struct {
	Handler
	ergo          *Ergo
	route         *Route
	method        string
	name          string
	description   string
	middleware    []Middleware
	params        map[string]*validation.Param
	schemes       []string
	consumes      []string
	produces      []string
	containsFiles bool
	bodyParams    bool
}

func NewOperation(handler Handler) *Operation {
	o := &Operation{
		params: map[string]*validation.Param{},
	}
	return o.SetHandler(handler)
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
	if method == constants.METHOD_GET ||
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

// Use allows adding middleware for o
// Use is additive
func (o *Operation) Use(middleware ...Middleware) *Operation {
	o.middleware = append(o.middleware, middleware...)
	return o
}

func (o *Operation) UseFunc(middleware ...MiddlewareFunc) *Operation {
	for _, m := range middleware {
		o.Use(Middleware(m))
	}
	return o
}

func (o *Operation) GetMiddleware() []Middleware {
	return o.middleware
}

// SetMiddleware will reset middleware with the given middleware
func (o *Operation) SetMiddleware(middleware ...Middleware) *Operation {
	o.middleware = middleware
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

func (o *Operation) Params(params ...*validation.Param) *Operation {
	addParams(o, params)
	return o
}

func (o *Operation) GetParams() map[string]*validation.Param {
	return o.params
}

func (o *Operation) GetParamsSlice() []*validation.Param {
	var params []*validation.Param
	for _, p := range o.params {
		params = append(params, p)
	}
	return params
}

func (o *Operation) SetParamsSlice(params ...*validation.Param) *Operation {
	o.setParamsSlice(params...)
	return o
}

func (o *Operation) SetParams(params map[string]*validation.Param) *Operation {
	o.setParams(params)
	return o
}

func (o *Operation) IgnoreParams(params ...string) *Operation {
	ignoreParams(o, params)
	return o
}

func (o *Operation) IgnoreParamsBut(params ...string) *Operation {
	ignoreParamsBut(o, params)
	return o
}

func (o *Operation) SetHandler(handler Handler) *Operation {
	if handler == nil {
		panic("Handler cannot be nil")
	}
	o.Handler = handler
	return o
}

func (o *Operation) Copy() *Operation {
	op := NewOperation(o.Handler)
	op.method = o.method
	op.description = o.description
	op.middleware = o.middleware
	op.params = o.params
	op.schemes = o.schemes
	op.consumes = o.consumes
	op.produces = o.consumes
	op.containsFiles = o.containsFiles
	return o
}

func (o *Operation) Validate(handler Handler) Handler {
	return HandlerFunc(func(ctx *Context) {

		// check for all things

		// not necessary as it's already been checked for in Route
		// if ctx.Request.Method != o.method {
		// 	//
		// 	return
		// }

		// check scheme
		var schemeAccepted bool
		if ctx.Request.URL.Scheme == "" {
			ctx.Request.URL.Scheme = constants.SCHEME_HTTP
		}

		schemeAccepted = containsString(o.schemes, ctx.Request.URL.Scheme)
		if !schemeAccepted {
			return
		}

		// check accept header
		var encoding string
		for _, acceptMime := range strings.Split(ctx.Request.Header.Get(constants.HEADER_Accept), ",") {
			mime := strings.Trim(strings.Split(acceptMime, ";")[0], " ")
			if 0 == len(mime) || mime == "*/*" {
				if len(o.produces) == 0 {
					encoding = DefMimeType
					break
				} else {
					encoding = o.produces[0]
					break
				}
			} else {
				if containsString(o.produces, mime) {
					encoding = mime
					break
				}
			}
		}

		if len(encoding) == 0 {
			ctx.Response.WriteString(http.StatusNotAcceptable, "406: Not Acceptable")
			return
		}

		ctx.Response.Encoding = encoding

		ctx.Ergo = o.ergo
		ctx.Operation = o
		q := ctx.Request.URL.Query()
		h := ctx.Request.Header

		if o.containsFiles {
			ctx.Request.ParseMultipartForm(MaxMemory)
		} else if o.bodyParams {
			ctx.Request.ParseForm()
		}
		for _, p := range o.params {
			var pv *validation.Value
			if p.IsInPath {
				v, ok := ctx.Request.PathParams[p.GetName()]
				if !ok {
					return
				}
				pv = validation.NewValue(p.GetName(), v, "path", p.GetAs())
			} else if p.IsInQuery {
				v, ok := q[p.GetName()]
				if !ok {
					if p.IsRequired {
						return
					}
				} else {
					if !p.IsMultiple {
						pv = validation.NewValue(p.GetName(), v[0], "query", p.GetAs())
					} else {
						pv = validation.NewMultipleValue(p.GetName(), v, "query", p.GetAs())
					}
				}
			} else if p.IsInHeader {
				v, ok := h[p.GetName()]
				if !ok {
					if p.IsRequired {
						return
					}
				} else {
					if !p.IsMultiple {
						pv = validation.NewValue(p.GetName(), v[0], "header", p.GetAs())
					} else {
						pv = validation.NewMultipleValue(p.GetName(), v, "header", p.GetAs())
					}
				}
			} else if p.IsInBody { // decide what to do when content type is form-encoded
				if p.IsFile {
					_, ok := ctx.Request.MultipartForm.File[p.GetName()]
					if !ok {
						if p.IsRequired {
							return
						}
					} else {
						//pv = NewFileParamValue(p.name, v[0], "header")
					}
				} else if !p.IsFile && o.containsFiles {
					_, ok := ctx.Request.MultipartForm.Value[p.GetName()]
					if !ok {
						if p.IsRequired {
							return
						}
					} else {
						//pv = NewFileParamValue(p.name, v[0], "header")
					}
				} else {
					v, ok := ctx.Request.Form[p.GetName()]
					if !ok {
						if p.IsRequired {
							return
						}
					} else {
						if !p.IsMultiple {
							pv = validation.NewValue(p.GetName(), v[0], "body", p.GetAs())
						} else {
							pv = validation.NewMultipleValue(p.GetName(), v, "body", p.GetAs())
						}
					}
				}
			}
			if pv != nil {
				ctx.Request.Input[p.GetName()] = pv
			}
		}
		for _, p := range o.params {
			pv, ok := ctx.Request.Input[p.GetName()]
			if ok {
				err := p.Validate(pv, ctx.Request)
				if err != nil {
					return
				}
			}
		}
		handler.ServeHTTP(ctx)
	})
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

func (o *Operation) setParams(params map[string]*validation.Param) {
	if params == nil {
		params = make(map[string]*validation.Param)
	}
	o.params = params
}

func (o *Operation) setParamsSlice(params ...*validation.Param) {
	paramsMap := map[string]*validation.Param{}
	for _, p := range params {
		paramsMap[p.GetName()] = p
	}
	o.setParams(paramsMap)
}
