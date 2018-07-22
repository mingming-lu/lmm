package ui

import (
	"fmt"
	"io"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/http"
	"lmm/api/testing"
)

func TestNewBlogTag_201(tt *testing.T) {
	t := testing.NewTester(tt)
	blog := randomBlog()

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	tag := Tag{
		Name: "a tag name",
	}

	res := postBlogTag(blog.ID(), headers, testing.StructToRequestBody(tag))
	t.Is(http.StatusCreated, res.StatusCode())
}

func TestNewBlogTag_401_NoAuthorizationHeader(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)

	tag := Tag{
		Name: "a tag name",
	}

	res := postBlogTag(1, headers, testing.StructToRequestBody(tag))
	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestNewBlogTag_400_InvalidPostBody(tt *testing.T) {
	t := testing.NewTester(tt)
	blog := randomBlog()

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := postBlogTag(blog.ID(), headers, testing.StructToRequestBody("?"))
	t.Is(http.StatusBadRequest, res.StatusCode())
}

func TestNewBlogTag_400_InvalidTagName(tt *testing.T) {
	t := testing.NewTester(tt)
	blog := randomBlog()

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	tag := Tag{
		Name: "$tag = ?",
	}

	res := postBlogTag(blog.ID(), headers, testing.StructToRequestBody(tag))
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(res.Body(), domain.ErrInvalidTagName.Error()+"\n")
}

func TestNewBlogTag_404_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)
	blog, err := factory.NewBlog(user.ID(), "title", "text")
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	tag := Tag{
		Name: "tag",
	}

	res := postBlogTag(blog.ID(), headers, testing.StructToRequestBody(tag))
	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(res.Body(), domain.ErrNoSuchBlog.Error()+"\n")
}

func postBlogTag(blogID uint64, headers map[string]string, requestBody io.Reader) *testing.Response {
	uri := fmt.Sprintf("/v1/blog/%d/tags", blogID)

	request := testing.POST(uri, requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.POST("/v1/blog/:blog/tags", accountUI.BearerAuth(ui.NewBlogTag))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
