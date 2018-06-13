package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type PathParams = httprouter.Params

type Path struct {
	raw    string
	params PathParams
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

func NewRequest(req *http.Request, params PathParams) *Request {
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

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
}

type responseWriter struct {
	headerWritten bool
	statusCode    int
	http.ResponseWriter
}

func newResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		headerWritten:  false,
		statusCode:     StatusOK,
		ResponseWriter: w,
	}
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if statusCode > 0 && rw.statusCode != statusCode {
		if rw.headerWritten {
			log.Printf("status code has been written as %d, cannot be written as %d again\n", rw.statusCode, statusCode)
			return
		}
		rw.statusCode = statusCode
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.headerWritten = true
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
	rw      ResponseWriter
	Request *Request
	values  Values
	logger  Logger
}

func (c *Context) Values() Values {
	return c.values
}

func (r *Request) ScanBody(schema interface{}) error {
	return json.NewDecoder(r.Request.Body).Decode(schema)
}

func NewContext(rw ResponseWriter, r *Request) *Context {
	return &Context{
		Request: r,
		rw:      rw,
		values:  make(Values),
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
