package main

import (
	"context"
	goHttp "net/http"
	"os"

	"lmm/api/pkg/http/middleware"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

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

	router := gin.New()
	router.Use(middleware.WrapAppEngineContext)

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

	goHttp.Handle("/", router)
	appengine.Main()
}
