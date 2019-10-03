package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"lmm/api/pkg/http/middleware"
	"lmm/api/pkg/pubsub"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine"

	// user
	userApp "lmm/api/service/user/application"
	userMessaging "lmm/api/service/user/port/adapter/messaging"
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
	dsClient     *datastore.Client
	gsClient     *storage.Client
	pubsubClient *pubsub.Client
)

var assetBucketName = os.Getenv("ASSET_BUCKET_NAME")

func initialze(c context.Context) func() {
	eg, egCtx := errgroup.WithContext(c)
	eg.Go(func() (err error) {
		gsClient, err = storage.NewClient(egCtx)
		return err
	})
	eg.Go(func() (err error) {
		dsClient, err = datastore.NewClient(egCtx, os.Getenv("DATASTORE_PROJECT_ID"))
		return err
	})
	eg.Go(func() (err error) {
		pubsubClient, err = pubsub.NewClient(c, os.Getenv("PUBSUB_PROJECT_ID"))
		return err
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}

	return func() {
		dsClient.Close()
		gsClient.Close()
		pubsubClient.Close()
	}
}

func main() {
	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	close := initialze(initCtx)
	defer close()

	// user
	userRepo := userStorage.NewUserDataStore(dsClient)
	userPub := userMessaging.NewUserEventPublisher(pubsubClient)
	userAppService := userApp.NewService(
		&userUtil.BcryptService{},
		userUtil.NewCFBTokenService(os.Getenv("LMM_API_TOKEN_KEY"), 24*time.Hour),
		userRepo,
		userRepo,
		userPub,
	)
	userUI := userUI.NewGinRouterProvider(userAppService)

	// article
	articleRepo := articleStorage.NewArticleDataStore(dsClient)
	articleUI := articleUI.NewGinRouterProvider(articleRepo, articleRepo, articleRepo)

	// asset
	assetRepo, err := assetStore.NewAssetDataStore(initCtx, dsClient, gsClient.Bucket(assetBucketName))
	if err != nil {
		panic(err)
	}
	assetStorage, err := assetStore.NewGCSUploader(initCtx, gsClient.Bucket(assetBucketName))
	if err != nil {
		panic(err)
	}
	assetUsecase := assetApp.New(assetRepo, assetStorage, assetRepo)
	assetUI := assetUI.NewGinRouterProvider(assetUsecase)

	router := gin.New()
	router.Use(middleware.CORS(os.Getenv("LMM_DOMAIN")), userUI.BearerAuth)

	userUI.Provide(router)
	articleUI.Provide(router)
	assetUI.Provide(router)

	http.Handle("/", router)
	appengine.Main()
}
