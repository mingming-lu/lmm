package ui

import (
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/infra"
	"lmm/api/http"
	"lmm/api/testing"
)

func TestGetAllTags_Success(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())
	testing.InitTable("tag")

	tags := make([]*model.Tag, 3)

	tag1, err := factory.NewTag(1, "b")
	t.NoError(err)
	t.NoError(repo.Add(tag1))
	tags[1] = tag1

	tag2, err := factory.NewTag(1, "c")
	t.NoError(err)
	t.NoError(repo.Add(tag2))
	tags[2] = tag2

	tag3, err := factory.NewTag(2, "a")
	t.NoError(err)
	t.NoError(repo.Add(tag3))
	tags[0] = tag3

	res := getAllTags()
	t.Is(http.StatusOK, res.StatusCode())
	t.JSON(tagsToJSON(tags), res.Body())
}

func TestGetAllTags_Empty(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	t := testing.NewTester(tt)
	testing.InitTable("tag")

	tags := make([]*model.Tag, 0)

	res := getAllTags()
	t.Is(http.StatusOK, res.StatusCode())
	t.JSON(tagsToJSON(tags), res.Body())
}

func getAllTags() *testing.Response {
	request := testing.GET("/v1/tags")

	router := testing.NewRouter()
	router.GET("/v1/tags", ui.GetAllTags)

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
