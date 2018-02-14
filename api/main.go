package main

import (
	"log"
	"net/http"

	"github.com/akinaru-lu/elesion"

	"lmm/api/blog"
	"lmm/api/db"
	"lmm/api/image"
	"lmm/api/profile"

	"lmm/api/controller/user"
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
	// router.GET("/blog", blog.GetBlog)
	// router.POST("/blog", blog.NewBlog)
	// router.PUT("/blog", blog.UpdateBlog)
	// router.DELETE("/blog", blog.DeleteBlog)

	// category
	// router.GET("/categories", blog.GetCategories)
	// router.POST("/categories", blog.NewCategory)

	// tag
	router.GET("/tags", blog.GetTags)

	router.GET("/photos", image.Handler)

	router.GET("/users/:user/profile", profile.GetProfile)

	log.Fatal(http.ListenAndServe(":8081", router))
}
