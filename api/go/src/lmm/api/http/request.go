package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

// Request is a abstraction of http request
type Request struct {
	pathParams  httprouter.Params
	queryParams url.Values
	*http.Request
}

func NewRequest(r *http.Request, params httprouter.Params) *Request {
	return &Request{pathParams: params, queryParams: nil, Request: r}
}

func (r *Request) Bind(schema interface{}) error {
	return json.NewDecoder(r.Request.Body).Decode(schema)
}

func (r *Request) PathParam(name string) string {
	return r.pathParams.ByName(name)
}

func (r *Request) QueryParam(name string) string {
	if r.queryParams == nil {
		r.queryParams = r.Request.URL.Query()
	}
	return r.queryParams.Get(name)
}

func (r *Request) RequestID() string {
	return r.Header.Get("X-Request-ID")
}

// RemoteAddr returns remote addr by cheking the following order:
// X-Real-IP header
// X-Forwarded-For header
// RemoteAddr property
func (r *Request) RemoteAddr() string {
	if remoteAddr := r.Header.Get("X-Real-IP"); remoteAddr != "" {
		return remoteAddr
	}

	if remoteAddr := r.Header.Get("X-Forwarded-For"); remoteAddr != "" {
		return remoteAddr
	}

	return r.Request.RemoteAddr
}

// Host returns host name by cheking the following order:
// Host header
// Host property
func (r *Request) Host() string {
	if host := r.Header.Get("Host"); host != "" {
		return host
	}

	return r.Request.Host
}

func (r *Request) Origin() string {
	return r.Header.Get("Origin")
}

func (r *Request) Path() string {
	return r.Request.URL.Path
}
