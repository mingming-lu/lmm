package ui

import (
	"fmt"
	"io"
	accountFactory "lmm/api/context/account/domain/factory"
	accountService "lmm/api/context/account/domain/service"
	accountStorage "lmm/api/context/account/infra"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/infra"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestUpdateBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := infra.NewBlogStorage(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	newTitle, newText := uuid.New(), uuid.New()
	requestBody := appservice.BlogContent{
		Title: newTitle,
		Text:  newText,
	}

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusOK, res.StatusCode())

	blog, err = repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(newTitle, blog.Title())
	t.Is(newText, blog.Text())
}

func TestUpdateBlog_Unauthorized(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := infra.NewBlogStorage(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := appservice.BlogContent{
		Title: "blog title",
		Text:  "blog text",
	}

	res := putBlog(nil, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestUpdateBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)

	requestBody := appservice.BlogContent{
		Title: "blog title",
		Text:  "blog text",
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(domain.ErrNoSuchBlog.Error()+"\n", res.Body())
}

func TestUpdateBlog_NoChange(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := infra.NewBlogStorage(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := appservice.BlogContent{
		Title: blog.Title(),
		Text:  blog.Text(),
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusNoContent, res.StatusCode())
}

func TestUpdateBlog_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := infra.NewBlogStorage(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := appservice.BlogContent{
		Title: "",
		Text:  blog.Text(),
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrEmptyBlogTitle.Error()+"\n", res.Body())
}

func TestUpdateBlog_NoPermission(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := infra.NewBlogStorage(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := appservice.BlogContent{
		Title: blog.Title(),
		Text:  blog.Text(),
	}

	otherUser, err := accountFactory.NewUser(uuid.New()[:31], uuid.New())
	t.NoError(err)
	accountRepo := accountStorage.NewUserStorage(testing.DB())
	t.NoError(accountRepo.Add(otherUser))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + accountService.EncodeToken(otherUser.Token())

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusForbidden, res.StatusCode())
	t.Is(domain.ErrNoSuchBlog.Error()+"\n", res.Body())
}

func putBlog(headers map[string]string, blogID uint64, requestBody io.Reader) *testing.Response {
	request := testing.PUT("/v1/blog/"+fmt.Sprint(blogID), requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.PUT("/v1/blog/:blog", accountUI.BearerAuth(ui.UpdateBlog))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
