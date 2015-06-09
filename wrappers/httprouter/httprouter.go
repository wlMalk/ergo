package httprouter

import (
	"net/http"

	"github.com/wlMalk/ergo"
	"github.com/wlMalk/ergo/wrappers"

	"github.com/julienschmidt/httprouter"
)

type Wrapper struct {
	router *httprouter.Router
}

func Wrap(router *httprouter.Router) *Wrapper {
	return &Wrapper{router}
}

func (w *Wrapper) Set(ops []*ergo.Operation) {
	for _, o := range ops {
		w.router.Handle(o.GetMethod(), wrappers.CurlyToColon(o.GetRoute().GetFullPath()), getHandle(o))
	}
}

func (w *Wrapper) Match(r *http.Request) http.Handler {
	h, ps, _ := w.router.Lookup(r.Method, r.URL.Path)
	if h == nil {
		return nil
	}
	return getHandler(h, ps)
}

func getHandle(h ergo.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		req := ergo.NewRequest(r)
		res := ergo.NewResponse(w)
		for _, p := range ps {
			req.PathParams[p.Key] = p.Value
		}
		h.ServeHTTP(ergo.NewContext(res, req))
	})
}

func getHandler(h httprouter.Handle, ps httprouter.Params) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, ps)
	})
}
