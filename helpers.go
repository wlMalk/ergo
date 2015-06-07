package ergo

import (
	"strings"

	"github.com/wlMalk/ergo/constants"
)

func preparePath(path string) string {
	path = strings.Trim(path, "/")
	return path
}

type paramer interface {
	GetParams() map[string]*Param
	setParams(map[string]*Param)
	GetParamsSlice() []*Param
	setParamsSlice(...*Param)
}

// addParams is additive, meaning that it will keep adding
// params as long as they are different in names.
// No two params can share the same name even if they are
// in different places.
func addParams(pa paramer, params ...*Param) {
	ps := pa.GetParamsSlice()
	for _, p := range params {
		ps = append(ps, p)
	}
	pa.setParamsSlice(ps...)
}

func ignoreParams(pa paramer, params ...string) {
	ps := pa.GetParams()
	for _, p := range params {
		delete(ps, p)
	}
	pa.setParams(ps)
}

func ignoreParamsBut(pa paramer, params ...string) {
	nparams := map[string]*Param{}
	ps := pa.GetParams()
	for _, p := range params {
		n, ok := ps[p]
		if ok {
			nparams[p] = n
		}
	}
	pa.setParams(nparams)
}

func prepareArgsSlice(args []string, f func(s string) bool) []string {
	if len(args) == 0 {
		return args
	}
	var nArgs []string
	for _, a := range args {
		a = strings.ToLower(a)
		duplicate := false
		for _, b := range nArgs {
			if a == b {
				duplicate = true
			}
		}
		if !duplicate && f(a) {
			nArgs = append(nArgs, a)
		}
	}

	return nArgs
}

type schemer interface {
	GetSchemes() []string
	setSchemes([]string)
}

func schemes(s schemer, schemes []string) {
	schemes = prepareArgsSlice(schemes, func(scheme string) bool {
		if scheme == constants.SCHEME_HTTP ||
			scheme == constants.SCHEME_HTTPS {
			return true
		}
		return false
	})
	if len(schemes) > 0 {
		s.setSchemes(schemes)
	}
}

type consumer interface {
	GetConsumes() []string
	setConsumes([]string)
}

func consumes(c consumer, mimes []string) {
	mimes = prepareArgsSlice(mimes, func(mime string) bool {
		if mime == constants.MIME_JSON ||
			mime == constants.MIME_XML {
			return true
		}
		return false
	})
	if len(mimes) > 0 {
		c.setConsumes(mimes)
	}
}

type producer interface {
	GetProduces() []string
	setProduces([]string)
}

func produces(p producer, mimes []string) {
	mimes = prepareArgsSlice(mimes, func(mime string) bool {
		if mime == constants.MIME_JSON ||
			mime == constants.MIME_XML {
			return true
		}
		return false
	})
	if len(mimes) > 0 {
		p.setProduces(mimes)
	}
}

func setOperation(r *Route, o *Operation) {
	o.setSchemes(r.ergo.GetSchemes())
	o.setConsumes(r.ergo.GetConsumes())
	o.setProduces(r.ergo.GetProduces())
	setParamer(r, o)
}

func setParamer(r *Route, p paramer) {
	p.setParams(r.GetParams())
}

func containsString(vals []string, a string) bool {
	for _, v := range vals {
		if a == v {
			return true
		}
	}
	return false
}

func containsInt(vals []int, a int) bool {
	for _, v := range vals {
		if a == v {
			return true
		}
	}
	return false
}
