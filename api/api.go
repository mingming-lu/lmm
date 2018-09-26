package api

import (
	accountInfra "lmm/api/context/account/infra"
	account "lmm/api/context/account/ui"
	articleFetcher "lmm/api/context/article/infra/fetcher"
	articlePersistence "lmm/api/context/article/infra/persistence"
	articleService "lmm/api/context/article/infra/service"
	article "lmm/api/context/article/ui"
	imageInfra "lmm/api/context/image/infra"
	img "lmm/api/context/image/ui"
	"lmm/api/http"
	"lmm/api/storage"
	"lmm/api/storage/static"
)

func NewRouter(db *storage.DB, cache *storage.Cache) *http.Router {
	userRepo := accountInfra.NewUserStorage(db)
	accountUI := account.New(userRepo)

	imgRepo := imageInfra.NewImageStorage(db)
	imgRepo.SetStaticRepository(static.NewLocalStaticRepository())
	imageUI := img.New(imgRepo)

	authorService := articleService.NewAuthorAdapter(db)
	articleFinder := articleFetcher.NewArticleFetcher(db)
	articleRepository := articlePersistence.NewArticleStorage(db, authorService)
	authorAdapter := articleService.NewAuthorAdapter(db)
	articleUI := article.NewUI(articleFinder, articleRepository, authorAdapter)

	router := http.NewRouter()

	// account
	router.POST("/v1/signup", accountUI.SignUp)
	router.POST("/v1/signin", accountUI.SignIn)
	router.GET("/v1/verify", accountUI.BearerAuth(accountUI.Verify))

	// article
	router.POST("/v1/articles", accountUI.BearerAuth(articleUI.PostArticle))
	router.PUT("/v1/articles/:articleID", accountUI.BearerAuth(articleUI.EditArticleText))
	router.GET("/v1/articles", articleUI.ListArticles)
	router.GET("/v1/articles/:articleID", articleUI.GetArticle)
	router.GET("/v1/articleTags", articleUI.GetAllArticleTags)

	// image
	router.POST("/v1/images", accountUI.BearerAuth(imageUI.Upload))
	router.GET("/v1/images", imageUI.LoadImagesByPage)
	router.PUT("/v1/images/:image", accountUI.BearerAuth(imageUI.MarkImage))

	return router
}
