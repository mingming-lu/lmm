package main

// TODO
// - versioning api

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
	router.GET("/blog/:blog", blog.Get)
	router.GET("/users/:user/blog", blog.GetList)
	router.POST("/blog", blog.Post)
	router.PUT("/blog/:blog", blog.Update)
	router.DELETE("/blog/:blog", blog.Delete)
	// blog category
	router.GET("/blog/:blog/category", blog.GetCategory)
	router.PUT("/blog/:blog/category", blog.SetCategory)
	router.DELETE("/blog/:blog/category", blog.DeleteCategory)

	// category
	router.GET("/users/:user/categories", category.GetByUser)
	router.POST("/categories", category.Register)
	router.PUT("/categories/:category", category.Update)

	// tag
	router.GET("/users/:user/tags", tag.GetByUser)
	router.GET("/blog/:blog/tags", tag.GetByBlog)
	router.POST("/blog/:blog/tags", tag.Register)
	router.PUT("/blog/:blog/tags/:tag", tag.Update)
	router.DELETE("/blog/:blog/tags/:tag", tag.Delete)

	// image
	router.GET("/users/:user/photos", image.GetPhotos)
	router.POST("/images", image.Upload)

	log.Fatal(http.ListenAndServe(":8081", router))
}
