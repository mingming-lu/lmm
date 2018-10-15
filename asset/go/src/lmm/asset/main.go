package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("/static")))

	log.Println("listening and serving at :8003")
	if err := http.ListenAndServe(":8003", nil); err != nil {
		log.Fatal(err.Error())
	}
}
