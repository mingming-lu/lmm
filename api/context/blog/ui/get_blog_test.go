package ui

import (
	"encoding/json"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestGetAllBlog_OK(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	repo := repository.NewBlogRepository()
	app := appservice.NewBlogApp(repo)

	t := testing.NewTester(tt)

	title, text := uuid.New(), uuid.New()
	app.PostNewBlog(user.ID(), title, text)

	res := getAllBlog()
	t.Is(http.StatusOK, res.StatusCode())

	blogListRes := BlogListResponse{}
	json.Unmarshal([]byte(res.Body()), &blogListRes)

	t.Is(1, len(blogListRes.Blog))
	t.Is(1, blogListRes.Page)
	t.False(blogListRes.HasNextPage)
}

func TestGetAllBlog_Empty(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)

	res := getAllBlog()
	t.Is(http.StatusOK, res.StatusCode())

	blogListRes := BlogListResponse{}
	json.Unmarshal([]byte(res.Body()), &blogListRes)

	t.Is(0, len(blogListRes.Blog))
	t.Is(1, blogListRes.Page)
	t.False(blogListRes.HasNextPage)
}

func getAllBlog() *testing.Response {
	request := testing.GET("/v1/blog")

	router := testing.NewRouter()
	router.GET("/v1/blog", GetAllBlog)

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
