package ui

import (
	"encoding/json"
	"fmt"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestGetBlog_OK(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	res := getBlog(blog.ID())

	t.Is(http.StatusOK, res.StatusCode())

	blogRes := BlogResponse{}
	json.Unmarshal([]byte(res.Body()), &blogRes)

	t.Is(fmt.Sprint(blog.ID()), blogRes.ID)
	t.Is(blog.Title(), blogRes.Title)
	t.Is(blog.Text(), blogRes.Text)
}

func TestGetBlog_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	// t.NoError(repo.Add(blog))

	res := getBlog(blog.ID())

	t.Is(http.StatusNotFound, res.StatusCode())
}

func getBlog(id uint64) *testing.Response {
	request := testing.GET("/v1/blog/" + fmt.Sprint(id))

	router := testing.NewRouter()
	router.GET("/v1/blog/:blog", GetBlog)

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
