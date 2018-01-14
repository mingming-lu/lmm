package article

import (
	"fmt"
	"os"
	"testing"

	"github.com/akinaru-lu/elesion"
	"github.com/stretchr/testify/assert"

	"lmm/api/db"
	"lmm/api/user"
)

var router *elesion.Router
var usr *user.User
var article *Article

func setUp() {
	db.Init("lmm_test")

	article, usr = NewTestArticle()

	router = elesion.New()
	router.GET("/articles", GetArticles)
	router.GET("/articles/:article", GetArticle)
	router.POST("/articles", NewArticle)
	router.PUT("/articles/:article", UpdateArticle)
	router.DELETE("/articles/:article", DeleteArticle)
}

func tearDown() {
	if err := db.New().DropDatabase("lmm_test").Close(); err != nil {
		fmt.Println(err)
	}
}

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()
	setUp()
	defer tearDown()
	code = m.Run()
}

func TestNewTestArticle(t *testing.T) {
	assert.Equal(t, int64(1), usr.ID)
	assert.Equal(t, int64(1), article.ID)
	assert.Equal(t, "test", article.Title)
	assert.Equal(t, "This is a test article", article.Text)
}
