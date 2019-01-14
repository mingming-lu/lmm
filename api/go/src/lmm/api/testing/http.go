package testing

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type RequestOptions struct {
	Headers map[string]string
}

func NewRequest(method, path string, body io.Reader, opts *RequestOptions) *http.Request {
	req := httptest.NewRequest(method, path, body)
	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Add(k, v)
		}
	}

	return req
}

func Do(req *http.Request, handler http.Handler) *Response {
	res := NewResponse()
	handler.ServeHTTP(res, req)

	return res
}

func GET(path string, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodGet, path, nil, opts)
}

func POST(path string, body io.Reader, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodPost, path, body, opts)
}

func PUT(path string, body io.Reader, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodPut, path, body, opts)
}

func DELETE(path string, opts *RequestOptions) *http.Request {
	return NewRequest(http.MethodDelete, path, nil, opts)
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
