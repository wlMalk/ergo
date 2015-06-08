package httprouter

import (
	"net/http"

	"github.com/wlMalk/ergo"

	"github.com/julienschmidt/httprouter"
)

type Wrapper struct {
	router *httprouter.Router
}

func Wrap(router *httprouter.Router) *Wrapper {
	return &Wrapper{router}
}

func (w *Wrapper) Handle(method string, path string, h ergo.Handler) {
	w.router.Handle(method, path, getHandle(h))
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
		params := map[string]string{}
		for _, p := range ps {
			params[p.Key] = p.Value
		}
		req.SetPathParams(params)
		h.ServeHTTP(res, req)
	})
}

func getHandler(h httprouter.Handle, ps httprouter.Params) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, ps)
	})
}
