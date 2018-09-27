package ui

import (
	"fmt"
	"io"
	"strings"

	"lmm/api/http"
	"lmm/api/service/article/domain"
	"lmm/api/testing"
)

func TestPostArticles_201(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := strings.Repeat("t", 140)
	body := "test body"
	res := postArticles(headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusCreated, res.StatusCode())
	t.Is("Success", res.Body())
	t.Regexp(`^\/v1\/articles\/\w+$`, res.Header().Get("Location"))
}

func TestPostArticles_401(tt *testing.T) {
	t := testing.NewTester(tt)

	res := postArticles(nil, testing.StructToRequestBody(struct{}{}))

	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestPostArticles_400_TitleRequired(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	dummy := ""
	res := postArticles(headers, testing.StructToRequestBody(struct {
		Body *string
		Tags []string
	}{
		Body: &dummy,
		Tags: make([]string, 0),
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(errTitleRequired.Error(), res.Body())
}

func TestPostArticles_400_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := ""
	body := "test body"
	res := postArticles(headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrEmptyArticleTitle.Error(), res.Body())
}

func TestPostArticles_400_InvalidTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := "$&%"
	body := "test body"
	res := postArticles(headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrInvalidArticleTitle.Error(), res.Body())
}

func TestPostArticles_400_TitleTooLong(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := strings.Repeat("t", 141)
	body := "test body"
	res := postArticles(headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrArticleTitleTooLong.Error(), res.Body())
}

func TestPostArticles_400_BodyRequired(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	dummy := ""
	res := postArticles(headers, testing.StructToRequestBody(struct {
		Title *string
		Tags  []string
	}{
		Title: &dummy,
		Tags:  []string{"test"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(errBodyRequired.Error(), res.Body())
}

func TestPostArticles_400_TagsRequired(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	dummy := ""
	res := postArticles(headers, testing.StructToRequestBody(struct {
		Title *string
		Body  *string
	}{
		Title: &dummy,
		Body:  &dummy,
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(errTagsRequired.Error(), res.Body())
}

func TestPostArticles_201_EmptyTags(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := "awesome"
	body := ""
	res := postArticles(headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{},
	}))

	t.Is(http.StatusCreated, res.StatusCode())
	t.Is("Success", res.Body())
}

func postArticles(headers map[string]string, requestBody io.ReadCloser) *testing.Response {
	uri := fmt.Sprint("/v1/articles")

	request := testing.POST(uri, requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
