package ui

import (
	"fmt"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/infra"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestDeleteTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := deleteTag(headers, tag.ID())
	t.Is(http.StatusNoContent, res.StatusCode())
}

func TestDeleteTag_Unauthorized(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)

	res := deleteTag(headers, tag.ID())
	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestDeleteTag_NoSuchTag(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := deleteTag(headers, tag.ID())
	t.Is(http.StatusNoContent, res.StatusCode())
}

func deleteTag(headers map[string]string, tagID uint64) *testing.Response {
	request := testing.DELETE("/v1/tags/" + fmt.Sprint(tagID))
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.DELETE("/v1/tags/:tag", accountUI.BearerAuth(ui.DeleteTag))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
