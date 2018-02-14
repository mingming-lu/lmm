package main

import (
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/controller/blog"
	"lmm/api/controller/user"
	"lmm/api/db"
	"lmm/api/image"
	"lmm/api/profile"
)

func init() {
	db.Init("lmm")
}

func main() {
	router := elesion.Default("[api]")

	// user
	router.POST("/signup", user.SignUp)
	router.POST("/signin", user.SignIn)
	router.GET("/verify", user.Verify)

	// blog
	router.GET("/blogs/:blog", blog.Get)
	router.GET("/users/:user/blogs", blog.GetByUser)
	router.POST("/blogs", blog.Post)
	router.PUT("/blogs/:blog", blog.Update)
	router.DELETE("/blogs/:blog", blog.Delete)

	// category
	// router.GET("/categories", blog.GetCategories)
	// router.POST("/categories", blog.NewCategory)

	// tag
	// router.GET("/tags", blog.GetTags)

	router.GET("/photos", image.Handler)

	router.GET("/users/:user/profile", profile.GetProfile)

	log.Fatal(http.ListenAndServe(":8081", router))
}
