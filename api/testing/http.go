package testing

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type Request struct {
	*http.Request
}

func NewRequest(method, path string, body io.Reader) *Request {
	return &Request{Request: httptest.NewRequest(method, path, body)}
}

func GET(path string) *Request {
	return NewRequest(http.MethodGet, path, nil)
}

func POST(path string, body io.Reader) *Request {
	return NewRequest(http.MethodPost, path, body)
}

func PUT(path string, body io.Reader) *Request {
	return NewRequest(http.MethodPut, path, body)
}

func DELETE(path string, body io.Reader) *Request {
	return NewRequest(http.MethodDelete, path, body)
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

func (r *Response) Body() string {
	b, err := ioutil.ReadAll(r.ResponseRecorder.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}
