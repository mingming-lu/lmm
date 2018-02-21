package testing

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func Request(method, path string, body io.Reader) *http.Request {
	return httptest.NewRequest(method, path, body)
}

func GET(path string) *http.Request {
	return Request(http.MethodGet, path, nil)
}

func POST(path string, body io.Reader) *http.Request {
	return Request(http.MethodPost, path, body)
}

func PUT(path string, body io.Reader) *http.Request {
	return Request(http.MethodPut, path, body)
}

func DELETE(path string, body io.Reader) *http.Request {
	return Request(http.MethodDelete, path, body)
}
