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

	// user
	router.POST("/signup", user.SignUp)
	router.POST("/login", user.Login)
	router.GET("/logout", user.Logout)
	router.GET("/verify", user.Verify)

	// article
	router.GET("/articles", article.GetArticles)
	router.POST("/articles", article.NewArticle)
	router.PUT("/articles", article.UpdateArticle)
	router.DELETE("/articles", article.DeleteArticle)

	// category
	router.GET("/categories", article.GetCategories)

	// tag
	router.GET("/tags", article.GetTags)

	router.GET("/photos", image.Handler)

	router.GET("/users/:user/profile", profile.GetProfile)

	log.Fatal(http.ListenAndServe(":8081", router))
}
