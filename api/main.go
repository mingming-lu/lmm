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
	router.GET("/users/:user/blogs", blog.GetList)
	router.POST("/blogs", blog.Post)
	router.PUT("/blogs/:blog", blog.Update)
	router.DELETE("/blogs/:blog", blog.Delete)
	// blog category
	router.GET("/blogs/:blog/category", blog.GetCategory)
	router.PUT("/blogs/:blog/category", blog.SetCategory)
	router.DELETE("/blogs/:blog/category", blog.DeleteCategory)

	// category
	router.GET("/users/:user/categories", category.GetByUser)
	router.POST("/categories", category.Register)
	router.PUT("/categories/:category", category.Update)

	// tag
	router.GET("/users/:user/tags", tag.GetByUser)
	router.GET("/blogs/:blog/tags", tag.GetByBlog)
	router.POST("/blogs/:blog/tags", tag.Register)
	router.PUT("/blogs/:blog/tags/:tag", tag.Update)
	router.DELETE("/blogs/:blog/tags/:tag", tag.Delete)

	// image
	router.GET("/users/:user/photos", image.GetPhotos)
	router.POST("/images", image.Upload)

	log.Fatal(http.ListenAndServe(":8081", router))
}
