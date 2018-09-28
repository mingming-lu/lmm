package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"lmm/api/http"
	"lmm/api/log"
	"lmm/api/middleware"
	accountInfra "lmm/api/service/account/infra"
	account "lmm/api/service/account/ui"
	articleFetcher "lmm/api/service/article/infra/fetcher"
	articlePersistence "lmm/api/service/article/infra/persistence"
	articleService "lmm/api/service/article/infra/service"
	article "lmm/api/service/article/ui"
	imageInfra "lmm/api/service/image/infra"
	img "lmm/api/service/image/ui"
	"lmm/api/storage"
	"lmm/api/storage/static"
)

func main() {
	logger := globalRecorder()
	defer logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	db := storage.NewDB()
	defer db.Close()

	_ = storage.NewCacheEngine()

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

	accessRecorder := accessLogRecorder()
	defer accessRecorder.Sync()
	router.Use(middleware.NewAccessLog(accessRecorder))

	recvRecoder := recoveryRecorder()
	defer recvRecoder.Sync()
	router.Use(middleware.NewRecovery(recvRecoder))

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

	http.Serve(":8002", router)
}

func globalRecorder() *zap.Logger {
	cfg := log.DefaultZapConfig()

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Named("global")
}

func recoveryRecorder() *zap.Logger {
	cfg := log.DefaultZapConfig()
	cfg.DisableCaller = true
	cfg.DisableStacktrace = false
	cfg.EncoderConfig.StacktraceKey = "stacktrace"

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Named("recovery")
}

func accessLogRecorder() *zap.Logger {
	cfg := log.DefaultZapConfig()
	cfg.DisableCaller = true

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Named("access_log")
}
