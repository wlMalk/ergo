package ergo

import (
	"strings"
)

func preparePath(path string) string {
	path = strings.Trim(path, "/")
	if path != "" {
		path = "/" + strings.ToLower(path)
	}
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
		if scheme == SCHEME_HTTP ||
			scheme == SCHEME_HTTPS {
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
		if mime == MIME_JSON ||
			mime == MIME_XML {
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
		if mime == MIME_JSON ||
			mime == MIME_XML {
			return true
		}
		return false
	})
	if len(mimes) > 0 {
		p.setProduces(mimes)
	}
}

type childer interface {
	paramer
	schemer
	consumer
	producer
}

func setChild(r *Route, child childer) {
	child.setParams(r.GetParams())
	child.setSchemes(r.GetSchemes())
	child.setConsumes(r.GetConsumes())
	child.setProduces(r.GetProduces())
}

