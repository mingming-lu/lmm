package ui

import (
	"context"
	"io"
	"net/http"

	"lmm/api/context/article/domain"
	"lmm/api/testing"
	"lmm/api/utils/strings"
)

func TestPutArticles_204(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := strings.Repeat("t", 140)
	body := "test body"
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusNoContent, res.StatusCode())
}

func TestPutArticles_401(tt *testing.T) {
	t := testing.NewTester(tt)

	res := putArticles("dummy123", nil, testing.StructToRequestBody(struct{}{}))

	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestPutArticles_403(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + testing.NewUser().Token()

	title := strings.Repeat("t", 140)
	body := "test body"
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusForbidden, res.StatusCode())
	t.Is(domain.ErrNotArticleAuthor.Error(), res.Body())
}

func TestPutArticles_404_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := strings.Repeat("t", 140)
	body := "test body"
	res := putArticles("88888888", headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(domain.ErrNoSuchArticle.Error(), res.Body())
}

func TestPutArticles_404_InvalidArticleID(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := strings.Repeat("t", 140)
	body := "test body"
	res := putArticles("dummy", headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(domain.ErrNoSuchArticle.Error(), res.Body())
}

func TestPutArticles_400_TitleRequired(tt *testing.T) {
	t := testing.NewTester(tt)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	dummy := ""
	res := putArticles("dummy123", headers, testing.StructToRequestBody(struct {
		Body *string
		Tags []string
	}{
		Body: &dummy,
		Tags: make([]string, 0),
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(errTitleRequired.Error(), res.Body())
}

func TestPutArticles_400_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := ""
	body := "test body"
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test", "testing"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrEmptyArticleTitle.Error(), res.Body())
}

func TestPutArticles_400_InvalidTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := "$&%"
	body := "test body"
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrInvalidArticleTitle.Error(), res.Body())
}

func TestPutArticles_400_TitleTooLong(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := strings.Repeat("t", 141)
	body := "test body"
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{"test"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrArticleTitleTooLong.Error(), res.Body())
}

func TestPutArticles_400_BodyRequired(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	dummy := ""
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(struct {
		Title *string
		Tags  []string
	}{
		Title: &dummy,
		Tags:  []string{"test"},
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(errBodyRequired.Error(), res.Body())
}

func TestPutArticles_400_TagsRequired(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	dummy := ""
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(struct {
		Title *string
		Body  *string
	}{
		Title: &dummy,
		Body:  &dummy,
	}))

	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(errTagsRequired.Error(), res.Body())
}

func TestPutArticles_204_EmptyTags(tt *testing.T) {
	t := testing.NewTester(tt)

	articleID, err := ui.appService.ArticleCommandService().PostNewArticle(context.Background(), user.ID(), "title", "body", []string{})
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	title := "awesome"
	body := ""
	res := putArticles(articleID.String(), headers, testing.StructToRequestBody(postArticleAdapter{
		Title: &title,
		Body:  &body,
		Tags:  []string{},
	}))

	t.Is(http.StatusNoContent, res.StatusCode())
}

func putArticles(articleID string, headers map[string]string, requestBody io.ReadCloser) *testing.Response {
	request := testing.PUT("/v1/articles/"+articleID, requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
