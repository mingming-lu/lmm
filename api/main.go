package main

// TODO
// - versioning api
// - change blogs to blog

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
	router.POST("/v1/signup", user.SignUp)
	router.POST("/v1/signin", user.SignIn)
	router.GET("/v1/verify", user.Verify)

	// blog
	router.GET("/v1/blogs/:blog", blog.Get)
	router.GET("/v1/users/:user/blogs", blog.GetList)
	router.POST("/v1/blogs", blog.Post)
	router.PUT("/v1/blogs/:blog", blog.Update)
	router.DELETE("/v1/blogs/:blog", blog.Delete)
	// blog category
	router.GET("/v1/blogs/:blog/category", blog.GetCategory)
	router.PUT("/v1/blogs/:blog/category", blog.SetCategory)
	router.DELETE("/v1/blogs/:blog/category", blog.DeleteCategory)

	// category
	router.GET("/v1/users/:user/categories", category.GetByUser)
	router.POST("/v1/categories", category.Register)
	router.PUT("/v1/categories/:category", category.Update)

	// tag
	router.GET("/v1/users/:user/tags", tag.GetByUser)
	router.GET("/v1/blogs/:blog/tags", tag.GetByBlog)
	router.POST("/v1/blogs/:blog/tags", tag.Register)
	router.PUT("/v1/blogs/:blog/tags/:tag", tag.Update)
	router.DELETE("/v1/blogs/:blog/tags/:tag", tag.Delete)

	// image
	router.GET("/v1/users/:user/photos", image.GetPhotos)
	router.POST("/v1/images", image.Upload)

	log.Fatal(http.ListenAndServe(":8081", router))
}
