package httptest

import (
	"io/ioutil"
	"net/http/httptest"
	"net/http"
	"io"
)

type ResponseWriter struct {
	*httptest.ResponseRecorder
}

func (w *ResponseWriter) Body() string {
	b, err := ioutil.ReadAll(w.ResponseRecorder.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (w *ResponseWriter) StatusCode() int {
	return w.Result().StatusCode
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{httptest.NewRecorder()}
}

func GET(path string) *http.Request {
	return httptest.NewRequest(http.MethodGet, path, nil)
}

func POST(path string, body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodPost, path, body)
}
