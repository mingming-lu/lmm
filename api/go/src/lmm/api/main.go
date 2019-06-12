package main

import (
	"context"
	"net/http"
	"os"

	"lmm/api/pkg/http/middleware"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	// user
	userApp "lmm/api/service/user/application"
	userStorage "lmm/api/service/user/port/adapter/persistence"
	userUI "lmm/api/service/user/port/adapter/presentation"
	userUtil "lmm/api/service/user/port/adapter/service"

	// article
	articleStorage "lmm/api/service/article/port/adapter/persistence"
	articleUI "lmm/api/service/article/port/adapter/presentation"

	// asset
	assetStore "lmm/api/service/asset/infra/persistence"
	assetUI "lmm/api/service/asset/presentation"
	assetApp "lmm/api/service/asset/usecase"
)

func main() {
	gcsClient, err := storage.NewClient(context.TODO())
	if err != nil {
		panic(err)
	}
	defer gcsClient.Close()

	datastoreClient, err := datastore.NewClient(context.TODO(), os.Getenv("DATASTORE_PROJECT_ID"))

	if err != nil {
		panic(err)
	}
	defer datastoreClient.Close()

	// user
	userRepo := userStorage.NewUserDataStore(datastoreClient)
	userAppService := userApp.NewService(&userUtil.BcryptService{}, &userUtil.CFBTokenService{}, userRepo, userRepo)
	userUI := userUI.NewGinRouterProvider(userAppService)

	// article
	articleRepo := articleStorage.NewArticleDataStore(datastoreClient)
	articleUI := articleUI.NewGinRouterProvider(articleRepo, articleRepo, articleRepo)

	// asset
	assetRepo := assetStore.NewAssetDataStore(datastoreClient)
	assetStorage := assetStore.NewGCSUploader(gcsClient)
	assetUsecase := assetApp.New(assetRepo, assetStorage, assetRepo)
	assetUI := assetUI.New(assetUsecase)

	router := gin.New()
	router.Use(middleware.WrapAppEngineContext, middleware.CORS(
		os.Getenv("APP_ORIGIN"),
		os.Getenv("MANAGER_ORIGIN"),
	))
	router.Use(userUI.BearerAuth(func(c *gin.Context) {
		// TODO: delete this
	}))

	router.POST("/v1/photos", userUI.BearerAuth(assetUI.PostV1Photos))
	router.GET("/v1/photos", assetUI.GetV1Photos)

	userUI.Provide(router)
	articleUI.Provide(router)

	http.Handle("/", router)
	appengine.Main()
}
