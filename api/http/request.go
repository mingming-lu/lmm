package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

// Request is a abstraction of http request
type Request interface {
	Bind(i interface{}) error

	Body() io.ReadCloser

	Header(name string) string

	PathParam(name string) string

	QueryParam(name string) string
}

type requestImpl struct {
	httprouter.Params
	*http.Request
	url.Values
}

func (r *requestImpl) Bind(schema interface{}) error {
	return json.NewDecoder(r.Body()).Decode(schema)
}

func (r *requestImpl) Body() io.ReadCloser {
	return r.Request.Body
}

func (r *requestImpl) Header(name string) string {
	return r.Request.Header.Get(name)
}

func (r *requestImpl) PathParam(name string) string {
	return r.Params.ByName(name)
}

func (r *requestImpl) QueryParam(name string) string {
	if r.Values == nil {
		r.Values = r.Request.URL.Query()
	}
	return r.Values.Get(name)
}
