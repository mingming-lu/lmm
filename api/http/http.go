package http

import (
	"log"
	"net/http"
)

type Handler = func(*Context)

func Serve(addr string, r *Router) {
	log.Fatal(http.ListenAndServe(addr, r))
}
