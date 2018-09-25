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
