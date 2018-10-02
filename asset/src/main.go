package main

import (
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func main() {
	http.DefaultServeMux.HandleFunc("/", index)

	if err := http.ListenAndServe(":8003", nil); err != nil {
		log.Fatal(err.Error())
	}
}
