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

	// /articles
	router.GET("/articles/:userID", articles.GetArticles)

	// /article
	router.GET("/article/:id", articles.GetArticle)
	router.POST("/article", articles.PostArticle)
	router.PUT("/article/:id", articles.UpdateArticle)
	// router.DELETE("/article/:id", articles.DeleteArticle)

	// /articles/categories
	router.GET("/articles/:userID/categories", articles.GetCategories)

	// /articles/category
	router.POST("/articles/category", articles.NewCategory)
	router.PUT("/articles/category/:id", articles.UpdateCategory)
	router.DELETE("/articles/category/:id", articles.DeleteCategory)

	// /articles/tags
	router.GET("/articles/:userID/tags", articles.GetTags)

	router.GET("/photos", image.Handler)

	router.GET("/profile", profile.Handler)

	log.Fatal(http.ListenAndServe(":8081", router))
}
