package main

import (
	"github.com/akinaru-lu/elesion"
	"lmm/api/article"
	"lmm/api/db"
	"lmm/api/user"
	"log"
	"net/http"
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

	// /article/categories
	// routing.GET("/article/:userID/categories", article.GetCategories)

	// /article/category
	// routing.POST("/article/category", article.NewCategory)
	// routing.PUT("/article/category/:id", article.UpdateCategory)
	// routing.DELETE("/article/category/:id", article.DeleteCategory)

	// /article/tags
	// routing.GET("/article/:userID/tags", article.GetTags)

	// /article/tags
	// routing.GET("/article/:id/tags", article.GetArticleTags)
	// routing.POST("/article/tags", article.NewTags)
	// routing.DELETE("/article/tags/:id", article.DeleteTag)

	// routing.GET("/photos", image.Handler)

	// routing.GET("/profile", profile.Handler)

	log.Fatal(http.ListenAndServe(":8081", router))
}
