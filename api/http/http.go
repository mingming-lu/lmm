package http

import (
	"log"
	"net/http"
)

const (
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
)

type Handler = func(*Context)

func Serve(addr string, r *Router) {
	log.Fatal(http.ListenAndServe(addr, r))
}

func HandleStatus(c *Context, code int) {
	c.Status(code).String(http.StatusText(code))
}

func BadRequest(c *Context) {
	HandleStatus(c, StatusBadRequest)
}

func Unauthorized(c *Context) {
	HandleStatus(c, StatusUnauthorized)
}

func NotFound(c *Context) {
	HandleStatus(c, StatusNotFound)
}

func InternalServerError(c *Context) {
	HandleStatus(c, StatusInternalServerError)
}

func Status(code int) string {
	return http.StatusText(code)
}
