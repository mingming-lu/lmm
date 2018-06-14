package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router *httprouter.Router
}

func NewRouter() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

func (r *Router) Handle(method string, path string, handler Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		c := NewContext(newResponseWriter(w), NewRequest(req, params))
		handler(c)
	})
}

func (r *Router) GET(path string, handler Handler) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.Handle(http.MethodDelete, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
