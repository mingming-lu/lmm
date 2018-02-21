package testing

import (
	"io/ioutil"
	"net/http/httptest"
)

type ResponseWriter struct {
	*httptest.ResponseRecorder
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{httptest.NewRecorder()}
}

func (w *ResponseWriter) StatusCode() int {
	return w.Result().StatusCode
}

func (w *ResponseWriter) Body() string {
	b, err := ioutil.ReadAll(w.ResponseRecorder.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}
