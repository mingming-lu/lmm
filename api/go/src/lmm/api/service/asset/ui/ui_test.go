package ui

import (
	"context"
	"net/http"
	"os"

	api "lmm/api/http"
	"lmm/api/service/asset/infra/persistence"
	"lmm/api/service/asset/infra/service"
	authApp "lmm/api/service/auth/application"
	authStorage "lmm/api/service/auth/infra/persistence"
	authUI "lmm/api/service/auth/ui"
	"lmm/api/storage/db"
	"lmm/api/testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	handler Handler
)

type NopUploader struct{}

func (u *NopUploader) Upload(c context.Context, id string, data []byte, args ...interface{}) error {
	return nil
}

func (u *NopUploader) Delete(c context.Context, id string, args ...interface{}) error {
	return nil
}

func (u *NopUploader) Close() error {
	return nil
}

type Handler func(*http.Request) *testing.Response

func NewHandler(db db.DB) Handler {
	repo := authStorage.NewUserStorage(db)
	app := authApp.NewService(repo)
	authUI := authUI.NewUI(app)

	assetUI := New(
		service.NewAssetFetcher(db),
		persistence.NewAssetStorage(db, &NopUploader{}),
		service.NewImageService(db),
		service.NewUserAdapter(db),
	)

	router := api.NewRouter()
	router.POST("/v1/assets/photos", authUI.BearerAuth(assetUI.UploadPhoto))
	router.PUT("/v1/assets/photos/:photo/alts", authUI.BearerAuth(assetUI.PutPhotoAlternateTexts))
	router.GET("/v1/assets/photos", assetUI.ListPhotos)

	return func(req *http.Request) *testing.Response {
		return testing.Do(req, router)
	}
}

func TestMain(m *testing.M) {
	mysql := db.DefaultMySQL()

	handler = NewHandler(mysql)

	code := m.Run()

	mysql.Close()
	os.Exit(code)
}
