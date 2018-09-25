package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"lmm/api/utils/strings"
)

var (
	timeout time.Duration
)

func init() {
	if s := os.Getenv("HTTP_TIMEOUT_SECOND"); s != "" {
		if i, err := strings.ParseUint(s); err != nil {
			timeout = time.Duration(i) * time.Second
		}
	}
}

// Router is used to handle routing
type Router struct {
	router *httprouter.Router
}

// NewRouter creates new router
func NewRouter() *Router {
	return &Router{
		router: httprouter.New(),
	}
}

// Handle registers handlers to handle combination of method and path
func (r *Router) Handle(method string, path string, handler Handler) {
	r.router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		var cancel context.CancelFunc

		if timeout > 0 {
			ctx, cancel = context.WithTimeout(r.Context(), timeout)
		}

		if cancel != nil {
			defer cancel()
		}

		c := &contextImpl{
			Context: ctx,
			r: &requestImpl{
				Params:  params,
				Request: r,
				Values:  nil,
			},
			w: w,
		}
		handler(c)
	})
}

// GET registers handler to handle GET method with given path
func (r *Router) GET(path string, handler Handler) {
	r.Handle(http.MethodGet, path, handler)
}

// POST registers handler to handle POST method with given path
func (r *Router) POST(path string, handler Handler) {
	r.Handle(http.MethodPost, path, handler)
}

// PUT registers handler to handle PUT method with given path
func (r *Router) PUT(path string, handler Handler) {
	r.Handle(http.MethodPut, path, handler)
}

// DELETE registers handler to handle DELETE method with given path
func (r *Router) DELETE(path string, handler Handler) {
	r.Handle(http.MethodDelete, path, handler)
}

// ServeHTTP implements http.Handler.ServeHTTP
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
