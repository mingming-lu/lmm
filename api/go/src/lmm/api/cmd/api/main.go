package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"lmm/api/http"
	"lmm/api/log"
	"lmm/api/messaging"
	"lmm/api/messaging/pubsub"
	"lmm/api/messaging/rabbitmq"
	"lmm/api/middleware"
	"lmm/api/storage/db"
	"lmm/api/storage/file"

	// user
	userApp "lmm/api/service/user/application"
	userEvent "lmm/api/service/user/domain/event"
	userMessaging "lmm/api/service/user/infra/messaging"
	userStorage "lmm/api/service/user/infra/persistence"
	userUI "lmm/api/service/user/ui"

	// auth
	authApp "lmm/api/service/auth/application"
	authStorage "lmm/api/service/auth/infra/persistence/mysql"
	authUI "lmm/api/service/auth/ui"

	// article
	articleFetcher "lmm/api/service/article/infra/fetcher"
	articleStorage "lmm/api/service/article/infra/persistence"
	authorService "lmm/api/service/article/infra/service"
	articleUI "lmm/api/service/article/ui"

	// asset
	assetDomainService "lmm/api/service/asset/domain/service"
	assetStorage "lmm/api/service/asset/infra/persistence"
	assetService "lmm/api/service/asset/infra/service"
	asset "lmm/api/service/asset/ui"
)

var (
	gcpProjectID string
)

func init() {
	gcpProjectID = os.Getenv("GCP_PROJECT_ID")
	if gcpProjectID == "" {
		panic("empty gcp project id")
	}
	fmt.Printf("gcp project id: %s\n", gcpProjectID)
}

func main() {
	pubsubClient, err := pubsub.NewClient(gcpProjectID, "/gcp/credentials/service_account.json")
	if err != nil {
		panic(err)
	}
	defer pubsubClient.Close()

	callback := log.Init(pubsub.NewPubSubTopicPublisher(
		pubsubClient.Topic(getEnvOrPanic("GCP_PUBSUB_TOPIC_API_LOG")),
		func() context.Context {
			return context.Background()
		},
	))
	defer callback()

	mysql := db.DefaultMySQL()
	defer mysql.Close()

	rabbitMQClient := rabbitmq.DefaultClient()
	rabbitMQUploader := file.NewRabbitMQAssetUploader(rabbitMQClient)
	defer rabbitMQUploader.Close() // would close rabbitMQClient too

	router := http.NewRouter()

	// middlewares
	// access log
	accessLogger := middleware.NewAccessLog(pubsub.NewPubSubTopicPublisher(
		pubsubClient.Topic(getEnvOrPanic("GCP_PUBSUB_TOPIC_API_ACCESS_LOG")),
		func() context.Context {
			return context.Background()
		},
	))
	defer accessLogger.Sync()
	router.Use(accessLogger.AccessLog)
	// recovery
	router.Use(middleware.Recovery)
	// request id
	router.Use(middleware.WithRequestID)

	// auth
	authRepo := authStorage.NewUserStorage(mysql)
	authAppService := authApp.NewService(authRepo)
	authUI := authUI.NewUI(authAppService)
	router.POST("/v1/auth/login", authUI.Login)

	// user
	userRepo := userStorage.NewUserStorage(mysql)
	userAppService := userApp.NewService(userRepo)
	userUI := userUI.NewUI(userAppService)
	userEventSubscriber := userMessaging.NewSubscriber(mysql)
	messaging.SyncBus().Subscribe(&userEvent.UserRoleChanged{}, userEventSubscriber.OnUserRoleChanged)
	messaging.SyncBus().Subscribe(&userEvent.UserPasswordChanged{}, userEventSubscriber.OnUserPasswordChanged)
	router.GET("/v1/users", authUI.BearerAuth(userUI.ViewAllUsers))
	router.POST("/v1/users", userUI.SignUp)
	router.PUT("/v1/users/:user/role", authUI.BearerAuth(userUI.AssignUserRole))
	router.PUT("/v1/users/:user/password", userUI.ChangeUserPassword)

	// article
	authorAdapter := authorService.NewAuthorAdapter(mysql)
	articleRepo := articleStorage.NewArticleStorage(mysql, authorAdapter)
	articleFinder := articleFetcher.NewArticleFetcher(mysql)
	articleUI := articleUI.NewUI(articleFinder, articleRepo, authorAdapter)
	router.POST("/v1/articles", authUI.BearerAuth(articleUI.PostNewArticle))
	router.PUT("/v1/articles/:articleID", authUI.BearerAuth(articleUI.EditArticle))
	router.GET("/v1/articles", articleUI.ListArticles)
	router.GET("/v1/articles/:articleID", articleUI.GetArticle)
	router.GET("/v1/articleTags", articleUI.GetAllArticleTags)

	// asset
	assetFinder := assetService.NewAssetFetcher(mysql)
	assetRepo := assetStorage.NewAssetStorage(mysql, rabbitMQUploader)
	imageService := assetService.NewImageService(mysql)
	imageEncoder := &assetDomainService.NopImageEncoder{}
	asset := asset.New(assetFinder, assetRepo, imageService, imageEncoder, assetService.NewUserAdapter(mysql))
	router.POST("/v1/assets/images", authUI.BearerAuth(asset.UploadImage))
	router.GET("/v1/assets/images", asset.ListImages)
	router.POST("/v1/assets/photos", authUI.BearerAuth(asset.UploadPhoto))
	router.PUT("/v1/assets/photos/:photo/alts", authUI.BearerAuth(asset.PutPhotoAlternateTexts))
	router.GET("/v1/assets/photos", asset.ListPhotos)
	router.GET("/v1/assets/photos/:photo", authUI.BearerAuth(asset.GetPhotoDescription))

	server := http.NewServer(":8002", router)
	server.Run()
}

func getEnvOrPanic(key string) string {
	s := os.Getenv(key)
	if s == "" {
		panic("empty environment variable: " + key)
	}
	return s
}
