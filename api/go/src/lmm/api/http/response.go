package http

import (
	"log"
	"net/http"
)

// Response interface used by Context
type Response interface {
	http.ResponseWriter
	StatusCode() int
}

type responseImpl struct {
	writer     http.ResponseWriter
	statusCode int
	written    bool
}

func newResponseImpl(rw http.ResponseWriter) *responseImpl {
	return &responseImpl{
		statusCode: http.StatusOK,
		writer:     rw,
		written:    false,
	}
}

func (r *responseImpl) Header() http.Header {
	return r.writer.Header()
}

func (r *responseImpl) StatusCode() int {
	return r.statusCode
}

func (r *responseImpl) Write(data []byte) (int, error) {
	return r.writer.Write(data)
}

func (r *responseImpl) WriteHeader(statusCode int) {
	if r.written {
		log.Print("unexpected to set status code more than once")
		return
	}
	r.writer.WriteHeader(statusCode)
	r.statusCode = statusCode
	r.written = true
}
