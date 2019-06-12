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
	assetStore "lmm/api/service/asset/port/adapter/persistence"
	assetUI "lmm/api/service/asset/port/adapter/presentation"
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
	assetUI := assetUI.NewGinRouterProvider(assetUsecase)

	router := gin.New()
	router.Use(middleware.WrapAppEngineContext, middleware.CORS(
		os.Getenv("APP_ORIGIN"),
		os.Getenv("MANAGER_ORIGIN"),
	), userUI.BearerAuth)

	userUI.Provide(router)
	articleUI.Provide(router)
	assetUI.Provide(router)

	http.Handle("/", router)
	appengine.Main()
}
