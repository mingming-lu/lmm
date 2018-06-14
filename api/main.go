package main

import (
	account "lmm/api/context/account/ui"
	blog "lmm/api/context/blog/ui"
	"lmm/api/usecase/auth"

	"lmm/api/http"
)

func main() {
	router := http.NewRouter()

	// account
	router.POST("/v1/signup", account.SignUp)
	router.POST("/v1/signin", account.SignIn)
	router.GET("/v1/verify", auth.BearerAuth(account.Verify))

	// // blog
	router.GET("/v1/blog", blog.GetAllBlog)
	router.GET("/v1/blog/:blog", blog.GetBlog)
	router.POST("/v1/blog", auth.BearerAuth(blog.PostBlog))
	router.PUT("/v1/blog/:blog", auth.BearerAuth(blog.UpdateBlog))
	// // blog category
	// router.GET("/v1/blog/:blog/category", blog.GetCategory)
	// router.PUT("/v1/blog/:blog/category", blog.SetCategory)
	// router.DELETE("/v1/blog/:blog/category", blog.DeleteCategory)

	// category
	router.POST("/v1/categories", blog.PostCategory)
	router.PUT("/v1/categories/:category", blog.UpdateCategory)
	// router.POST("/v1/categories", category.Register)
	// router.PUT("/v1/categories/:category", category.Update)
	// router.DELETE("/v1/categories/:category", category.Delete)

	// // tag
	// router.GET("/v1/users/:user/tags", tag.GetByUser)
	// router.GET("/v1/blog/:blog/tags", tag.GetByBlog)
	// router.POST("/v1/blog/:blog/tags", tag.Register)
	// router.PUT("/v1/blog/:blog/tags/:tag", tag.Update)
	// router.DELETE("/v1/blog/:blog/tags/:tag", tag.Delete)

	// // image
	// router.GET("/v1/users/:user/images", image.GetAllImages)
	// router.GET("/v1/users/:user/images/photos", image.GetPhotos)
	// router.POST("/v1/images", image.Upload)
	// router.PUT("/v1/images/putPhoto", image.PutPhoto)
	// router.PUT("/v1/images/removePhoto", image.RemovePhoto)

	// log.Fatal(http.ListenAndServe(":8002", router))
	http.Serve(":8002", router)
}
