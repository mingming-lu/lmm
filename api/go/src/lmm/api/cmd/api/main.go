package main

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/http"
	"lmm/api/log"
	"lmm/api/messaging/pubsub"
	"lmm/api/messaging/rabbitmq"
	"lmm/api/middleware"
	"lmm/api/storage/db"
	"lmm/api/storage/uploader"

	// user
	userApp "lmm/api/service/user/application"
	userStorage "lmm/api/service/user/infra/persistence"
	userUI "lmm/api/service/user/ui"

	// auth
	authApp "lmm/api/service/auth/application"
	authStorage "lmm/api/service/auth/infra/persistence"
	authUI "lmm/api/service/auth/ui"

	// article
	articleFetcher "lmm/api/service/article/infra/fetcher"
	articleStorage "lmm/api/service/article/infra/persistence"
	authorService "lmm/api/service/article/infra/service"
	articleUI "lmm/api/service/article/ui"

	// asset
	assetStorage "lmm/api/service/asset/infra/persistence"
	assetService "lmm/api/service/asset/infra/service"
	asset "lmm/api/service/asset/ui"
)

func main() {
	pubsubClient, err := pubsub.NewClient()
	if err != nil {
		panic(err)
	}
	defer pubsubClient.Close()

	pubsubTopicPublisher := pubsub.NewPubSubTopicPublisher(
		pubsubClient.Topic(os.Getenv("GCP_PROJECT_TOPIC_ID")),
		func() context.Context {
			return context.Background()
		},
	)

	callback := log.Init(pubsubTopicPublisher)
	defer callback()

	mysql := db.DefaultMySQL()
	defer mysql.Close()

	// localUploader := uploader.NewLocalImageUploader()
	rabbitMQClient := rabbitmq.DefaultClient()
	rabbitMQUploader := uploader.NewRabbitMQAssetUploader(rabbitMQClient)
	defer rabbitMQUploader.Close() // would close rabbitMQClient too

	router := http.NewRouter()

	// middlewares begin
	// access logg
	accessLogger := middleware.NewAccessLog(pubsubTopicPublisher)
	defer accessLogger.Sync()
	router.Use(accessLogger.AccessLog)
	// recovery
	router.Use(middleware.Recovery)
	// request id
	router.Use(middleware.WithRequestID)
	// cache control
	router.Use(middleware.CacheControl)

	// user
	userRepo := userStorage.NewUserStorage(mysql)
	userAppService := userApp.NewService(userRepo)
	userUI := userUI.NewUI(userAppService)
	router.POST("/v1/users", userUI.SignUp)

	// auth
	authRepo := authStorage.NewUserStorage(mysql)
	authAppService := authApp.NewService(authRepo)
	authUI := authUI.NewUI(authAppService)
	router.POST("/v1/auth/login", authUI.Login)

	// article
	authorAdapter := authorService.NewAuthorAdapter(mysql)
	articleRepo := articleStorage.NewArticleStorage(mysql, authorAdapter)
	articleFinder := articleFetcher.NewArticleFetcher(mysql)
	articleUI := articleUI.NewUI(articleFinder, articleRepo, authorAdapter)
	router.POST("/v1/articles", authUI.BearerAuth(articleUI.PostNewArticle))
	router.PUT("/v1/articles/:articleID", authUI.BearerAuth(articleUI.EditArticle))
	router.GET("/v1/articles", articleUI.ListArticles)
	router.GET("/v2/articles", articleUI.ListArticlesByPagination)
	router.GET("/v1/articles/:articleID", articleUI.GetArticle)
	router.GET("/v1/articleTags", articleUI.GetAllArticleTags)

	// asset
	assetFinder := assetService.NewAssetFetcher(mysql)
	assetRepo := assetStorage.NewAssetStorage(mysql, rabbitMQUploader)
	asset := asset.New(assetFinder, assetRepo, assetService.NewUserAdapter(mysql))
	router.POST("/v1/assets/images", authUI.BearerAuth(asset.UploadImage))
	router.GET("/v1/assets/images", asset.ListImages)
	router.POST("/v1/assets/photos", authUI.BearerAuth(asset.UploadPhoto))
	router.GET("/v1/assets/photos", asset.ListPhotos)

	server := http.NewServer(":8002", router)
	server.Run()
}
