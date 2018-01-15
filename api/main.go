package main

import (
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/article"
	"lmm/api/db"
	"lmm/api/image"
	"lmm/api/profile"
	"lmm/api/user"
)

func init() {
	db.Init("lmm")
}

func main() {
	router := elesion.Default("[api]")

	// users
	router.GET("/users/:user", user.GetUser)
	router.POST("/users", user.NewUser)

	// article
	router.GET("/users/:user/articles", article.GetArticles)
	router.GET("/users/:user/articles/:article", article.GetArticle)
	router.POST("/users/:user/articles", article.NewArticle)
	router.PUT("/users/:user/articles/:article", article.UpdateArticle)
	router.DELETE("/users/:user/articles/:article", article.DeleteArticle)

	// /users/:user/categories
	router.GET("/users/:user/categories", article.GetCategories)

	// /users/:user/tags
	router.GET("/users/:user/tags", article.GetTags)
	router.GET("/users/:user/articles/:article/tags", article.GetArticleTags)

	router.GET("/photos", image.Handler)

	router.GET("/users/:user/profile", profile.GetProfile)

	log.Fatal(http.ListenAndServe(":8081", router))
}
