package ui

import (
	"io"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain/service"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestPostBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	content := appservice.BlogContent{
		Title: "blog title",
		Text:  "blog text",
	}

	res := postBlog(headers, testing.StructToRequestBody(content))
	t.Is(http.StatusCreated, res.StatusCode())
	t.Regexp(`/blog/\d+`, res.Header().Get("Location"))
}

func TestPostBlog_Unauthorized(tt *testing.T) {
	t := testing.NewTester(tt)

	content := appservice.BlogContent{
		Title: "blog title",
		Text:  "blog text",
	}

	res := postBlog(nil, testing.StructToRequestBody(content))
	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestPostBlog_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	blog := appservice.BlogContent{
		Title: "",
		Text:  uuid.New(),
	}

	res := postBlog(headers, testing.StructToRequestBody(blog))
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(service.ErrEmptyBlogTitle.Error()+"\n", res.Body())
}

func postBlog(headers map[string]string, requestBody io.Reader) *testing.Response {
	request := testing.POST("/v1/blog", requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.POST("/v1/blog", accountUI.BearerAuth(ui.PostBlog))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
