package ui

import (
	"encoding/json"
	"lmm/api/context/blog/appservice"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestGetAllBlog_OK(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)

	content := appservice.BlogContent{
		Title: uuid.New(),
		Text:  uuid.New(),
	}

	body := testing.StructToRequestBody(content)
	ui.app.PostNewBlog(user, body)

	res := getAllBlog()
	t.Is(http.StatusOK, res.StatusCode())

	blogListPage := appservice.BlogListPage{}
	json.Unmarshal([]byte(res.Body()), &blogListPage)

	t.Is(1, len(blogListPage.Blog))
	t.Is(-1, blogListPage.NextPage)
}

func TestGetAllBlog_Empty(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)

	res := getAllBlog()
	t.Is(http.StatusOK, res.StatusCode())

	blogListPage := appservice.BlogListPage{}
	json.Unmarshal([]byte(res.Body()), &blogListPage)

	t.Is(0, len(blogListPage.Blog))
	t.Is(-1, blogListPage.NextPage)
}

func getAllBlog() *testing.Response {
	request := testing.GET("/v1/blog")

	router := testing.NewRouter()
	router.GET("/v1/blog", ui.GetAllBlog)

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
