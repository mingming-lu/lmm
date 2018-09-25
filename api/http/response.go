package http

import "net/http"

type Response interface {
	http.ResponseWriter
}
