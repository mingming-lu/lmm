package main

import (
	accountInfra "lmm/api/context/account/infra"
	account "lmm/api/context/account/ui"
	blogInfra "lmm/api/context/blog/infra"
	blog "lmm/api/context/blog/ui"

	"lmm/api/http"
	"lmm/api/storage"
)

var (
	accountUI *account.UI
	blogUI    *blog.UI
)

func initUIs(db *storage.DB) {
	userRepo := accountInfra.NewUserStorage(db)
	accountUI = account.New(userRepo)

	blogRepo := blogInfra.NewBlogStorage(db)
	categoryRepo := blogInfra.NewCategoryStorage(db)
	tagRepo := blogInfra.NewTagStorage(db)
	blogUI = blog.New(blogRepo, categoryRepo, tagRepo)
}

func main() {
	db := storage.NewDB()
	initUIs(db)

	router := http.NewRouter()

	// account
	router.POST("/v1/signup", accountUI.SignUp)
	router.POST("/v1/signin", accountUI.SignIn)
	router.GET("/v1/verify", accountUI.BearerAuth(accountUI.Verify))

	// blog
	router.GET("/v1/blog", blogUI.GetAllBlog)
	router.GET("/v1/blog/:blog", blogUI.GetBlog)
	router.POST("/v1/blog", accountUI.BearerAuth(blogUI.PostBlog))
	router.PUT("/v1/blog/:blog", accountUI.BearerAuth(blogUI.UpdateBlog))
	// blog category
	router.GET("/v1/blog/:blog/category", blogUI.GetBlogCagetory)
	router.PUT("/v1/blog/:blog/category", blogUI.SetBlogCategory)

	// category
	router.GET("/v1/categories", blogUI.GetAllCategoris)
	router.POST("/v1/categories", accountUI.BearerAuth(blogUI.PostCategory))
	router.PUT("/v1/categories/:category", accountUI.BearerAuth(blogUI.UpdateCategory))
	router.DELETE("/v1/categories/:category", accountUI.BearerAuth(blogUI.DeleteCategory))

	http.Serve(":8002", router)
}

// tag
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
