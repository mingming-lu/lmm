package main

import (
	"log"
	"net/http"

	"./articles"
	"./image"
	"./profile"

    _ "github.com/astaxie/beego"
)

type Router map[string]func(w http.ResponseWriter, r *http.Request)

var router = Router{
	"/articles": articles.HandleArticles,
	"/photos":   image.HandlePhotos,
	"/profile":  profile.HandleProfile,
}

func main() {
	for path, handler := range router {
		http.HandleFunc(path, handler)
	}
	log.Fatal(http.ListenAndServe(":8081", nil))
}
