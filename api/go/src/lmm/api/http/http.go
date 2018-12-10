package http

import (
	"net/http"

	"go.uber.org/zap"
)

const (
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusNoContent           = http.StatusNoContent
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusRequestTimeout      = http.StatusRequestTimeout
	StatusConflict            = http.StatusConflict
	StatusInternalServerError = http.StatusInternalServerError
	StatusServiceUnavailable  = http.StatusServiceUnavailable
)

const (
	// self-defined status codes

	// StatusClientAbort defines the code when client aborted before response from server
	StatusClientAbort = 477
)

type Handler = func(Context)

type Middleware = func(Handler) Handler

func Serve(addr string, r *Router) {
	zap.L().Info("Serving at:" + addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}

func HandleStatus(c Context, code int) {
	c.String(code, StatusText(code))
}

func Created(c Context) {
	HandleStatus(c, http.StatusCreated)
}

func NoContent(c Context) {
	HandleStatus(c, http.StatusNoContent)
}

func BadRequest(c Context) {
	HandleStatus(c, StatusBadRequest)
}

func Unauthorized(c Context) {
	HandleStatus(c, StatusUnauthorized)
}

func NotFound(c Context) {
	HandleStatus(c, StatusNotFound)
}

func RequestTimeout(c Context) {
	HandleStatus(c, StatusRequestTimeout)
}

func ClientAbort(c Context) {
	HandleStatus(c, StatusClientAbort)
}

func InternalServerError(c Context) {
	HandleStatus(c, StatusInternalServerError)
}

func ServiceUnavailable(c Context) {
	HandleStatus(c, StatusServiceUnavailable)
}

func StatusText(code int) string {
	if code == StatusClientAbort {
		return "Client Abort"
	}

	return http.StatusText(code)
}
