package testing

import (
	"io/ioutil"
	"net/http/httptest"
)

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
