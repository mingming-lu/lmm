package ui

import (
	"context"
	"os"

	"lmm/api/http"
	userUI "lmm/api/service/user/ui"
	"lmm/api/testing"

	"cloud.google.com/go/datastore"
)

var (
	router *http.Router
)

func TestMain(m *testing.M) {
	c := context.Background()

	dataStore, err := datastore.NewClient(c, "")
	if err != nil {
		panic("failed to connect to datastore: " + err.Error())
	}

	userUI.NewUI(nil)

	code := m.Run()

	dataStore.Close()

	os.Exit(code)
}

func TestPostV1Articles(t testing.T) {
	router := http.NewRouter()
	router.POST("/v1/articles", userUI.BearerAuth())
}
