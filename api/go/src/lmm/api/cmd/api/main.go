package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"lmm/api/http"
	"lmm/api/log"
	"lmm/api/messaging/pubsub"
	"lmm/api/messaging/rabbitmq"
	"lmm/api/middleware"
	"lmm/api/storage/file"

	// user
	userApp "lmm/api/service/user/application"
	userStorage "lmm/api/service/user/infra/persistence"
	userUtil "lmm/api/service/user/infra/service"
	userUI "lmm/api/service/user/ui"

	// article
	articleStorage "lmm/api/service/article/infra/persistence"
	articleUI "lmm/api/service/article/ui"

	// asset
	assetStore "lmm/api/service/asset/infra/persistence"
	assetUI "lmm/api/service/asset/presentation"
	assetApp "lmm/api/service/asset/usecase"
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

	gcsClient, err := storage.NewClient(context.TODO(), option.WithCredentialsFile("/gcp/credentials/service_account.json"))
	if err != nil {
		panic(err)
	}
	defer gcsClient.Close()

	callback := log.Init(pubsub.NewPubSubTopicPublisher(
		pubsubClient.Topic(getEnvOrPanic("GCP_PUBSUB_TOPIC_API_LOG")),
		func() context.Context {
			return context.Background()
		},
	))
	defer callback()

	datastoreClient, err := datastore.NewClient(context.TODO(), "lmm")
	if err != nil {
		panic(err)
	}
	defer datastoreClient.Close()

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

	// user
	userRepo := userStorage.NewUserDataStore(datastoreClient)
	userAppService := userApp.NewService(&userUtil.BcryptService{}, &userUtil.CFBTokenService{}, userRepo, userRepo)
	userUI := userUI.NewUI(userAppService)
	router.POST("/v1/users", userUI.SignUp)
	router.PUT("/v1/users/:user/password", userUI.ChangeUserPassword)
	router.POST("/v1/auth/token", userUI.Token)

	// article
	articleRepo := articleStorage.NewArticleDataStore(datastoreClient)
	articleUI := articleUI.NewUI(articleRepo, articleRepo, articleRepo)
	router.POST("/v1/articles", userUI.BearerAuth(articleUI.PostNewArticle))
	router.PUT("/v1/articles/:articleID", userUI.BearerAuth(articleUI.PutV1Articles))
	router.GET("/v1/articles", articleUI.ListArticles)
	router.GET("/v1/articles/:articleID", articleUI.GetArticle)
	router.GET("/v1/articleTags", articleUI.GetAllArticleTags)

	// asset
	assetRepo := assetStore.NewAssetDataStore(datastoreClient)
	assetStorage := assetStore.NewGCSUploader(gcsClient)
	assetUsecase := assetApp.New(assetRepo, assetStorage, assetRepo)
	assetUI := assetUI.New(assetUsecase)

	router.POST("/v1/photos", userUI.BearerAuth(assetUI.PostV1Photos))
	router.GET("/v1/photos", assetUI.GetV1Photos)

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
