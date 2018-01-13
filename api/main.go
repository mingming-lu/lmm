package main

import (
	"lmm/api/article"
	"lmm/api/db"
	"log"
	"net/http"
	"github.com/akinaru-lu/elesion"
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

	// /article/categories
	// router.GET("/article/:userID/categories", article.GetCategories)

	// /article/category
	// router.POST("/article/category", article.NewCategory)
	// router.PUT("/article/category/:id", article.UpdateCategory)
	// router.DELETE("/article/category/:id", article.DeleteCategory)

	// /article/tags
	// router.GET("/article/:userID/tags", article.GetTags)

	// /article/tags
	// router.GET("/article/:id/tags", article.GetArticleTags)
	// router.POST("/article/tags", article.NewTags)
	// router.DELETE("/article/tags/:id", article.DeleteTag)

	// router.GET("/photos", image.Handler)

	// router.GET("/profile", profile.Handler)

	log.Fatal(http.ListenAndServe(":8081", router))
}
