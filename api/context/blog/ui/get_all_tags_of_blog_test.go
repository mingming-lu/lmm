package ui

import (
	"fmt"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/infra"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestGetAllTagsOfBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())
	tagRepo := infra.NewTagStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tags := make([]*model.Tag, 2)

	tag1, err := factory.NewTag(blog.ID(), "c")
	t.NoError(err)
	t.NoError(tagRepo.Add(tag1))
	tags[1] = tag1

	tag2, err := factory.NewTag(blog.ID(), "b")
	t.NoError(err)
	t.NoError(tagRepo.Add(tag2))
	tags[0] = tag2

	tag3, err := factory.NewTag(testing.GenerateID(), "a")
	t.NoError(err)
	t.NoError(tagRepo.Add(tag3))

	res := getAllTagsOfBlog(blog.ID())
	t.Is(http.StatusOK, res.StatusCode())
	t.JSON(tagsToJSON(tags), res.Body())
}

func TestGetAllTagsOfBlog_Empty(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tags := make([]*model.Tag, 0)

	res := getAllTagsOfBlog(blog.ID())
	t.Is(http.StatusOK, res.StatusCode())
	t.JSON(tagsToJSON(tags), res.Body())
}

func TestGetAllTagsOfBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)

	res := getAllTagsOfBlog(blog.ID())
	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(domain.ErrNoSuchBlog.Error()+"\n", res.Body())
}

func getAllTagsOfBlog(blogID uint64) *testing.Response {
	request := testing.GET(fmt.Sprintf("/v1/blog/%d/tags", blogID))

	router := testing.NewRouter()
	router.GET("/v1/blog/:blog/tags", ui.GetAllTagsOfBlog)

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
