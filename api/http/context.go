package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type URLParams = httprouter.Params

type Path struct {
	raw    string
	params URLParams
}

func (p *Path) String() string {
	return p.raw
}

func (p *Path) Params(name string) string {
	return p.params.ByName(name)
}

type URL struct {
	*url.URL
	Path  *Path
	query url.Values
}

func (u *URL) Query(name string) string {
	return u.query.Get(name)
}

type Request struct {
	*http.Request
	*URL
}

func NewRequest(req *http.Request, params URLParams) *Request {
	path := &Path{
		raw:    req.URL.Path,
		params: params,
	}
	url := &URL{
		URL:   req.URL,
		Path:  path,
		query: req.URL.Query(),
	}
	return &Request{
		Request: req,
		URL:     url,
	}
}

type Values map[string]interface{}

func (vs Values) Set(key string, v interface{}) {
	vs[key] = v
}

func (vs Values) Get(key string) interface{} {
	return vs[key]
}

type Context struct {
	rw      http.ResponseWriter
	Request *Request
	values  Values
}

func (c *Context) Values() Values {
	return c.values
}

func (r *Request) ScanBody(schema interface{}) error {
	return json.NewDecoder(r.Request.Body).Decode(schema)
}

func NewContext(rw http.ResponseWriter, r *http.Request, ps URLParams) *Context {
	return &Context{
		Request: NewRequest(r, ps),
		rw:      rw,
	}
}

func (c *Context) Status(code int) *Context {
	c.rw.WriteHeader(code)
	return c
}

func (c *Context) Header(key, value string) *Context {
	c.rw.Header().Add(key, value)
	return c
}

func (c *Context) JSON(data interface{}) *Context {
	c.rw.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(c.rw)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
	return c
}

func (c *Context) String(s string) *Context {
	c.rw.Header().Set("Content-Type", "text/plain")
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	c.rw.Write([]byte(s))
	return c
}

func (c *Context) Stringf(format string, a ...interface{}) *Context {
	return c.String(fmt.Sprintf(format, a...))
}
