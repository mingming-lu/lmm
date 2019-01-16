package ui

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"

	"lmm/api/http"
	"lmm/api/service/article/domain"
	"lmm/api/testing"
	"lmm/api/util/stringutil"
	"lmm/api/util/testutil"
)

func TestPostArticles(tt *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	t := testing.NewTester(tt)
	user := testutil.NewUser(mysql)

	cases := map[string]struct {
		ReqTitle      *string
		ReqBody       *string
		ReqTags       []string
		ReqHeaders    http.Header
		ResStatusCode int
		ResBody       string
		ResHeaders    map[string]string
	}{
		"Created": {
			ReqTitle: stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:  stringutil.Pointer(uuid.New().String()),
			ReqTags: func() []string {
				s := make([]string, 0)
				for i := 0; i < rand.Intn(10); i++ {
					s = append(s, uuid.New().String()[:8])
				}
				return s
			}(),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusCreated,
			ResBody:       "Success",
			ResHeaders:    map[string]string{"Location": `^\/v1\/articles\/\w+$`},
		},
		"NoTags": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusCreated,
			ResBody:       "Success",
			ResHeaders:    map[string]string{"Location": `^\/v1\/articles\/\w+$`},
		},
		"Unauthorized": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    nil,
			ResStatusCode: http.StatusUnauthorized,
			ResBody:       http.StatusText(http.StatusUnauthorized),
			ResHeaders:    nil,
		},
		"TitleRequired": {
			ReqTitle:      nil,
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTitleRequired.Error(),
			ResHeaders:    nil,
		},
		"BodyRequired": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       nil,
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errBodyRequired.Error(),
			ResHeaders:    nil,
		},
		"TagsRequired": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       nil,
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTagsRequired.Error(),
			ResHeaders:    nil,
		},
		"EmptyTitle": {
			ReqTitle:      stringutil.Pointer(""),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrEmptyArticleTitle.Error(),
			ResHeaders:    nil,
		},
		"InvalidTitle": {
			ReqTitle:      stringutil.Pointer("!@#$"),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrInvalidArticleTitle.Error(),
			ResHeaders:    nil,
		},
		"LongTitle": {
			ReqTitle:      stringutil.Pointer(strings.Repeat("t", 141)),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrArticleTitleTooLong.Error(),
			ResHeaders:    nil,
		},
	}

	for testName, testCase := range cases {
		t.Run(testName, func(_ *testing.T) {
			res := postArticles(&testing.RequestOptions{
				Headers: testCase.ReqHeaders,
				FormData: testing.StructToRequestBody(
					postArticleAdapter{
						Title: testCase.ReqTitle,
						Body:  testCase.ReqBody,
						Tags:  testCase.ReqTags,
					},
				),
			})
			t.Is(testCase.ResStatusCode, res.StatusCode())
			t.Is(testCase.ResBody, res.Body())
			for k, v := range testCase.ResHeaders {
				t.Regexp(v, res.Header().Get(k))
			}
		})
	}
}

func postArticles(opts *testing.RequestOptions) *testing.Response {
	return testing.DoRequest(testing.POST("/v1/articles", opts), router)
}
