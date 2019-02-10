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
	decoder := json.NewDecoder(r.Request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(schema)
}

func (r *Request) PathParam(name string) string {
	return r.pathParams.ByName(name)
}

func (r *Request) QueryParam(name string) *string {
	r.parseQuery()
	if ps := r.queryParams[name]; len(ps) != 0 {
		return &ps[0]
	}
	return nil
}

func (r *Request) QueryParamOrDefault(name string, v string) string {
	p := r.QueryParam(name)
	if p == nil {
		return v
	}
	return *p
}

func (r *Request) QueryParams(name string) []string {
	r.parseQuery()
	return r.queryParams[name]
}

func (r *Request) parseQuery() {
	if r.queryParams == nil {
		r.queryParams = r.Request.URL.Query()
	}
}

func (r *Request) RequestID() string {
	return r.Header.Get("X-Request-ID")
}

// ClientIP returns the client ip by cheking the following order:
// X-Real-IP header
// RemoteAddr property
func (r *Request) ClientIP() string {
	if remoteAddr := r.Header.Get("X-Real-IP"); remoteAddr != "" {
		return remoteAddr
	}

	return r.Request.RemoteAddr
}

// HostName returns host name by cheking the following order:
// Host header
// Host property
func (r *Request) HostName() string {
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
