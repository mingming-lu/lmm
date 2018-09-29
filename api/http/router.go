package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"lmm/api/util/stringutil"
)

var (
	timeout = 5 * time.Second
)

func init() {
	if s := os.Getenv("HTTP_TIMEOUT_SECOND"); s != "" {
		if i, err := stringutil.ParseUint(s); err == nil && i > 0 {
			timeout = time.Duration(i) * time.Second
		} else {
			zap.L().Warn(err.Error())
		}
	}
}

// Router is used to handle routing
type Router struct {
	middlewares []Middleware
	router      *httprouter.Router
}

// NewRouter creates new router
func NewRouter() *Router {
	return &Router{
		middlewares: make([]Middleware, 0),
		router:      httprouter.New(),
	}
}

// Handle registers handlers to handle combination of method and path
func (r *Router) Handle(method string, path string, handler Handler) {
	matryoshka := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		matryoshka = r.middlewares[i](matryoshka)
	}

	r.router.Handle(method, path, func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		c := &contextImpl{
			req: NewRequest(req.WithContext(ctx), params),
			res: newResponseImpl(res),
		}

		matryoshka(c)
	})
}

// Use registers a middleware to router
func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
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
