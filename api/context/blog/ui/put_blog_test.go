package ui

import (
	"fmt"
	"io"
	accountFactory "lmm/api/context/account/domain/factory"
	accountRepository "lmm/api/context/account/domain/repository"
	accountService "lmm/api/context/account/domain/service"
	"lmm/api/context/blog/appservice"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/usecase/auth"
	"lmm/api/utils/uuid"
)

func TestUpdateBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := repository.NewBlogRepository()

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	newTitle, newText := uuid.New(), uuid.New()
	requestBody := Blog{
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

	repo := repository.NewBlogRepository()

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := Blog{
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

	requestBody := Blog{
		Title: "blog title",
		Text:  "blog text",
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(appservice.ErrNoSuchBlog.Error()+"\n", res.Body())
}

func TestUpdateBlog_NoChange(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := repository.NewBlogRepository()

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := Blog{
		Title: blog.Title(),
		Text:  blog.Text(),
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusNoContent, res.StatusCode())
}

func TestUpdateBlog_NoPermission(tt *testing.T) {
	t := testing.NewTester(tt)

	repo := repository.NewBlogRepository()

	title, text := uuid.New(), uuid.New()
	blog, err := factory.NewBlog(user.ID(), title, text)
	t.NoError(err)
	t.NoError(repo.Add(blog))

	requestBody := Blog{
		Title: blog.Title(),
		Text:  blog.Text(),
	}

	otherUser, err := accountFactory.NewUser(uuid.New()[:31], uuid.New())
	t.NoError(err)
	accountRepo := accountRepository.New()
	t.NoError(accountRepo.Add(otherUser))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + accountService.EncodeToken(otherUser.Token())

	res := putBlog(headers, blog.ID(), testing.StructToRequestBody(requestBody))
	t.Is(http.StatusForbidden, res.StatusCode())
	t.Is(appservice.ErrNoSuchBlog.Error()+"\n", res.Body())
}

func putBlog(headers map[string]string, blogID uint64, requestBody io.Reader) *testing.Response {
	request := testing.PUT("/v1/blog/"+fmt.Sprint(blogID), requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.PUT("/v1/blog/:blog", auth.BearerAuth(UpdateBlog))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
