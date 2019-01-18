package ui

import (
	"regexp"
	"strings"

	"lmm/api/http"
	"lmm/api/service/article/domain"
	"lmm/api/testing"
	"lmm/api/util/stringutil"
	"lmm/api/util/testutil"

	"github.com/google/uuid"
)

func TestPutArticlews(tt *testing.T) {
	lock.Lock()
	defer lock.Unlock()

	user := testutil.NewUser(mysql)

	res := postArticles(
		&testing.RequestOptions{
			Headers: http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			FormData: testing.StructToRequestBody(postArticleAdapter{
				Title: stringutil.Pointer("title"),
				Body:  stringutil.Pointer("body"),
				Tags:  []string{"tag"},
			}),
		},
	)

	if res.StatusCode() != http.StatusCreated {
		tt.Fatal("failed to create test article data")
	}

	groups := regexp.MustCompile(`^/v1/articles/(\w+)$`).FindStringSubmatch(res.Header().Get("Location"))
	articleID := groups[1]

	cases := map[string]struct {
		ArticleID     string
		ReqTitle      *string
		ReqBody       *string
		ReqTags       []string
		ReqHeaders    http.Header
		ResStatusCode int
		ResBody       string
	}{
		"Success": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       []string{"foo", "bar"},
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusNoContent,
			ResBody:       "",
		},
		"NoTags": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusNoContent,
			ResBody:       "",
		},
		"Unauthorized": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    nil,
			ResStatusCode: http.StatusUnauthorized,
			ResBody:       http.StatusText(http.StatusUnauthorized),
		},
		"Forbidden": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + testutil.NewUser(mysql).AccessToken()}},
			ResStatusCode: http.StatusForbidden,
			ResBody:       domain.ErrNotArticleAuthor.Error(),
		},
		"NotFound": {
			ArticleID:     "notfound",
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + testutil.NewUser(mysql).AccessToken()}},
			ResStatusCode: http.StatusNotFound,
			ResBody:       domain.ErrNoSuchArticle.Error(),
		},
		"InvalidArticleID": {
			ArticleID:     "over8charcters",
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + testutil.NewUser(mysql).AccessToken()}},
			ResStatusCode: http.StatusNotFound,
			ResBody:       domain.ErrNoSuchArticle.Error(),
		},
		"TitleRequired": {
			ArticleID:     articleID,
			ReqTitle:      nil,
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTitleRequired.Error(),
		},
		"BodyRequired": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       nil,
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errBodyRequired.Error(),
		},
		"TagsRequired": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       nil,
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTagsRequired.Error(),
		},
		"EmptyTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(""),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrEmptyArticleTitle.Error(),
		},
		"InvalidTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("!@#$"),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrInvalidArticleTitle.Error(),
		},
		"LongTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(strings.Repeat("t", 141)),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrArticleTitleTooLong.Error(),
		},
	}

	for testName, testCase := range cases {
		tt.Run(testName, func(tt *testing.T) {
			t := testing.NewTester(tt)
			res := putArticles(testCase.ArticleID, &testing.RequestOptions{
				Headers: testCase.ReqHeaders,
				FormData: testing.StructToRequestBody(postArticleAdapter{
					Title: testCase.ReqTitle,
					Body:  testCase.ReqBody,
					Tags:  testCase.ReqTags,
				},
				),
			})
			t.Is(testCase.ResStatusCode, res.StatusCode())
			t.Is(testCase.ResBody, res.Body())
		})
	}
}

func putArticles(articleID string, opts *testing.RequestOptions) *testing.Response {
	request := testing.PUT("/v1/articles/"+articleID, opts)

	return testing.DoRequest(request, router)
}
