package persistence

import (
	"os"

	"lmm/api/storage"
	"lmm/api/testing"
)

var (
	articleStorage *ArticleStorage
)

func TestMain(m *testing.M) {
	db := storage.NewDB()
	defer db.CloseNow()

	articleStorage = NewArticleStorage(db)

	code := m.Run()
	os.Exit(code)
}
