package ui

import (
	"io"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/usecase/auth"
)

func TestPostBlog_Success(t *testing.T) {
	tester := testing.NewTester(t)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	blog := Blog{
		Title: "blog title",
		Text:  "blog text",
	}

	res := postBlog(headers, testing.StructToRequestBody(blog))
	tester.Is(http.StatusCreated, res.StatusCode())
	tester.Regexp(`/blog/\d+`, res.Header().Get("Location"))
}

func TestPostBlog_Unauthorized(t *testing.T) {
	tester := testing.NewTester(t)

	blog := Blog{
		Title: "blog title",
		Text:  "blog text",
	}

	res := postBlog(nil, testing.StructToRequestBody(blog))
	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func postBlog(headers map[string]string, requestBody io.Reader) *testing.Response {
	request := testing.POST("/v1/blog", requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.POST("/v1/blog", auth.BearerAuth(PostBlog))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
