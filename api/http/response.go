package http

import (
	"net/http"

	"go.uber.org/zap"
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
		zap.L().Warn("unexpected to set status code more than once",
			zap.Int("current", r.statusCode),
			zap.Int("input", statusCode),
		)
		return
	}
	r.writer.WriteHeader(statusCode)
	r.statusCode = statusCode
	r.written = true
}
