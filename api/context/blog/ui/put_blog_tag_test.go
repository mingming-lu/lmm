package ui

import (
	"fmt"
	"io"

	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/infra"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestUpdateTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	reqeustBody := Tag{
		Name: "tag",
	}

	res := putTag(tag.ID(), headers, testing.StructToRequestBody(reqeustBody))
	t.Is(http.StatusOK, res.StatusCode())
}

func TestUpdateTag_Unauthorized(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)

	reqeustBody := Tag{
		Name: "tag",
	}

	res := putTag(tag.ID(), headers, testing.StructToRequestBody(reqeustBody))
	t.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestUpdateTag_InvalidRequestBody(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	res := putTag(tag.ID(), headers, testing.StructToRequestBody("{aa}"))
	t.Is(http.StatusBadRequest, res.StatusCode())
}

func TestUpdateTag_InvalidTagName(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewTagStorage(testing.DB())

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	reqeustBody := Tag{
		Name: "",
	}

	res := putTag(tag.ID(), headers, testing.StructToRequestBody(reqeustBody))
	t.Is(http.StatusBadRequest, res.StatusCode())
	t.Is(domain.ErrInvalidTagName.Error()+"\n", res.Body())
}

func TestUpdateTag_NoSuchTag(tt *testing.T) {
	t := testing.NewTester(tt)

	tag, err := factory.NewTag(1, uuid.New()[:31])
	t.NoError(err)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + user.Token()

	reqeustBody := Tag{
		Name: "tag",
	}

	res := putTag(tag.ID(), headers, testing.StructToRequestBody(reqeustBody))
	t.Is(http.StatusNotFound, res.StatusCode())
	t.Is(domain.ErrNoSuchTag.Error()+"\n", res.Body())
}

func putTag(tagID uint64, headers map[string]string, requestBody io.Reader) *testing.Response {
	uri := fmt.Sprintf("/v1/tags/%d", tagID)

	request := testing.PUT(uri, requestBody)
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	router := testing.NewRouter()
	router.PUT("/v1/tags/:tag", accountUI.BearerAuth(ui.UpdateTag))

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
