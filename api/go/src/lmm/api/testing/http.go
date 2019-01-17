package testing

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// RequestOptions defines options used by request
type RequestOptions struct {
	Headers         http.Header
	FormData        io.Reader // should convert json or url-encoded string into io.Reader
	QueryParameters url.Values
}

// NewRequest creates a new request for testing
func NewRequest(method, path string, opts *RequestOptions) *http.Request {
	if opts == nil {
		opts = &RequestOptions{}
	}

	req := httptest.NewRequest(method, path, opts.FormData)

	// copy headers
	for head, values := range opts.Headers {
		for _, value := range values {
			req.Header.Add(head, value)
		}
	}

	// copy query parameters
	q := req.URL.Query()
	for name, values := range opts.QueryParameters {
		for _, value := range values {
			q.Add(name, value)
		}
	}

	return req
}

// DoRequest does request and return a response for testing
func DoRequest(req *http.Request, handler http.Handler) *Response {
	res := NewResponse()

	handler.ServeHTTP(res, req)

	return res
}

// GET creates a GET request
func GET(path string, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodGet, path, opts)
}

// POST creates a POST request
func POST(path string, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodPost, path, opts)
}

// PUT creates a PUT request
func PUT(path string, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodPut, path, opts)
}

// DELETE creates a DELETE request
func DELETE(path string, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodDelete, path, opts)
}

type Response struct {
	*httptest.ResponseRecorder
}

func NewResponse() *Response {
	return &Response{httptest.NewRecorder()}
}

func (r *Response) StatusCode() int {
	return r.Result().StatusCode
}

func (r *Response) RawBody() io.Reader {
	return r.ResponseRecorder.Body
}

func (r *Response) Body() string {
	b, err := ioutil.ReadAll(r.ResponseRecorder.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}
