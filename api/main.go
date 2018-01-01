package main

import (
	"lmm/api/articles"
	"lmm/api/db"
	"lmm/api/image"
	"lmm/api/profile"
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"
)

func init() {
	db.Init()
}

func main() {
	router := elesion.Default("[api]")
	router.GET("/articles", articles.GetArticles)
	router.GET("/article", articles.GetArticle)
	router.GET("/articles/categories", articles.GetCategories)
	router.GET("/photos", image.Handler)
	router.GET("/profile", profile.Handler)
	log.Fatal(http.ListenAndServe(":8081", router))
}
