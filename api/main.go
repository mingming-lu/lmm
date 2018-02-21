package main

import (
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/controller/blog"
	"lmm/api/controller/category"
	"lmm/api/controller/image"
	"lmm/api/controller/tag"
	"lmm/api/controller/user"
	"lmm/api/db"
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
	router.GET("/users/:user/categories", category.GetByUser)
	router.GET("/blogs/:blog/category", category.GetByBlog)
	router.POST("/blogs/:blog/category", category.Register)
	router.PUT("/blogs/:blog/category", category.Update)

	// tag
	router.GET("/users/:user/tags", tag.GetByUser)
	router.GET("/blogs/:blog/tags", tag.GetByBlog)
	router.POST("/blogs/:blog/tags", tag.Register)
	router.PUT("/blogs/:blog/tags/:tag", tag.Update)
	router.DELETE("/blogs/:blog/tags/:tag", tag.Delete)

	// image
	// router.GET("/users/:user/photos", image.Handler)
	router.POST("/images", image.Upload)

	router.GET("/users/:user/profile", profile.GetProfile)

	log.Fatal(http.ListenAndServe(":8081", router))
}
