package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"lmm/api/pkg/http/middleware"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

var (
	dsClient *datastore.Client
	gsClient *storage.Client
)

func init() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eg, egCtx := errgroup.WithContext(timeoutCtx)
	eg.Go(func() (err error) {
		gsClient, err = storage.NewClient(egCtx)
		return err
	})
	eg.Go(func() (err error) {
		dsClient, err = datastore.NewClient(egCtx, os.Getenv("DATASTORE_PROJECT_ID"))
		return err
	})

	if err := eg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%#v", err)
		os.Exit(1)
	}
}

func main() {
	defer dsClient.Close()
	defer gsClient.Close()

	// user
	userRepo := userStorage.NewUserDataStore(dsClient)
	userAppService := userApp.NewService(&userUtil.BcryptService{}, &userUtil.CFBTokenService{}, userRepo, userRepo)
	userUI := userUI.NewGinRouterProvider(userAppService)

	// article
	articleRepo := articleStorage.NewArticleDataStore(dsClient)
	articleUI := articleUI.NewGinRouterProvider(articleRepo, articleRepo, articleRepo)

	// asset
	assetRepo := assetStore.NewAssetDataStore(dsClient)
	assetStorage := assetStore.NewGCSUploader(gsClient)
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
