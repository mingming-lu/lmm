package main

import (
	account "lmm/api/context/account/ui"
	"lmm/api/storage"

	"lmm/api/http"
)

func main() {
	db := storage.NewDB()
	router := http.NewRouter()

	accountUI := account.New(db)
	// account
	router.POST("/v1/signup", accountUI.SignUp)
	router.POST("/v1/signin", accountUI.SignIn)
	router.GET("/v1/verify", accountUI.BearerAuth(accountUI.Verify))

	// blogUI := blog.New(db)
	// // // blog
	// router.GET("/v1/blog", blog.GetAllBlog)
	// router.GET("/v1/blog/:blog", blog.GetBlog)
	// router.POST("/v1/blog", auth.BearerAuth(blogUI.PostBlog))
	// router.PUT("/v1/blog/:blog", auth.BearerAuth(blog.UpdateBlog))
	// // // blog category
	// router.GET("/v1/blog/:blog/category", blog.GetBlogCagetory)
	// router.PUT("/v1/blog/:blog/category", blog.SetBlogCategory)

	// // category
	// router.GET("/v1/categories", blog.GetAllCategoris)
	// router.POST("/v1/categories", auth.BearerAuth(blog.PostCategory))
	// router.PUT("/v1/categories/:category", auth.BearerAuth(blog.UpdateCategory))
	// router.DELETE("/v1/categories/:category", auth.BearerAuth(blog.DeleteCategory))

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
