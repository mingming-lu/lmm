package article

import (
	"fmt"
	"testing"

	"github.com/akinaru-lu/elesion"

	"lmm/api/db"
)

var router = elesion.New()

func TestMain(m *testing.M) {
	router.GET("/articles", GetArticles)
	router.GET("/articles/:article", GetArticle)
	router.POST("/articles", NewArticle)
	router.PUT("/articles/:article", UpdateArticle)
	router.DELETE("/articles/:article", DeleteArticle)

	db.Init("lmm_test")
	defer func() {
		err := db.New().DropDatabase("lmm_test").Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	m.Run()
}
