package mux

import (
	"net/http"

	"github.com/wlMalk/ergo"

	"github.com/gorilla/mux"
)

type paramHandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

func (f paramHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r, nil)
}

type Wrapper struct {
	router *mux.Router
}

func Wrap(router *mux.Router) *Wrapper {
	return &Wrapper{router}
}

func (w *Wrapper) Set(routes []*ergo.Route) {
	for _, r := range routes {
		w.router.Handle(r.GetFullPath(), getHandle(r))
	}
}

func (w *Wrapper) Match(r *http.Request) http.Handler {
	rm := &mux.RouteMatch{}
	m := w.router.Match(r, rm)
	if !m {
		return nil
	}
	return getHandler(rm.Handler.(paramHandlerFunc), rm.Vars)
}

func getHandle(h ergo.Handler) http.Handler {
	return paramHandlerFunc(func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		req := ergo.NewRequest(r)
		res := ergo.NewResponse(w)
		req.PathParams = ps
		h.ServeHTTP(ergo.NewContext(res, req))
	})
}

func getHandler(h paramHandlerFunc, ps map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, ps)
	})
}
