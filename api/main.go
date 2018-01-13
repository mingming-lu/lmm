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
	router.POST("/users/articles", article.NewArticle)
	router.GET("/users/:user/articles/:article", article.GetArticle)

	// tag

	// /article
	// router.GET("/article/:id", article.GetArticle)
	// router.POST("/article", article.NewArticle)
	// router.PUT("/article/:id", article.UpdateArticle)
	// router.DELETE("/article/:id", article.DeleteArticle)

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
