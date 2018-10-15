package ui

import (
	"io"
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
	t := testing.NewTester(tt)
	user := testutil.NewUser(mysql)

	res := postArticles(
		map[string]string{"Authorization": "Bearer " + user.AccessToken()},
		testing.StructToRequestBody(postArticleAdapter{
			Title: stringutil.Pointer("title"),
			Body:  stringutil.Pointer("body"),
			Tags:  []string{"tag"},
		}),
	)

	if !t.Is(http.StatusCreated, res.StatusCode()) {
		t.FailNow()
	}

	groups := regexp.MustCompile(`^/v1/articles/(\w+)$`).FindStringSubmatch(res.Header().Get("Location"))
	articleID := groups[1]

	cases := map[string]struct {
		ArticleID     string
		ReqTitle      *string
		ReqBody       *string
		ReqTags       []string
		ReqHeaders    map[string]string
		ResStatusCode int
		ResBody       string
	}{
		"Success": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       []string{"foo", "bar"},
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusNoContent,
			ResBody:       http.StatusText(http.StatusNoContent),
		},
		"NoTags": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusNoContent,
			ResBody:       http.StatusText(http.StatusNoContent),
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
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + testutil.NewUser(mysql).AccessToken()},
			ResStatusCode: http.StatusForbidden,
			ResBody:       domain.ErrNotArticleAuthor.Error(),
		},
		"NotFound": {
			ArticleID:     "notfound",
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + testutil.NewUser(mysql).AccessToken()},
			ResStatusCode: http.StatusNotFound,
			ResBody:       domain.ErrNoSuchArticle.Error(),
		},
		"InvalidArticleID": {
			ArticleID:     "over8charcters",
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + testutil.NewUser(mysql).AccessToken()},
			ResStatusCode: http.StatusNotFound,
			ResBody:       domain.ErrNoSuchArticle.Error(),
		},
		"TitleRequired": {
			ArticleID:     articleID,
			ReqTitle:      nil,
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTitleRequired.Error(),
		},
		"BodyRequired": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       nil,
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errBodyRequired.Error(),
		},
		"TagsRequired": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       nil,
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTagsRequired.Error(),
		},
		"EmptyTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(""),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrEmptyArticleTitle.Error(),
		},
		"InvalidTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("!@#$"),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrInvalidArticleTitle.Error(),
		},
		"LongTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(strings.Repeat("t", 141)),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    map[string]string{"Authorization": "Bearer " + user.AccessToken()},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrArticleTitleTooLong.Error(),
		},
	}

	for testName, testCase := range cases {
		t.Run(testName, func(_ *testing.T) {
			res := putArticles(testCase.ArticleID, testCase.ReqHeaders, testing.StructToRequestBody(postArticleAdapter{
				Title: testCase.ReqTitle,
				Body:  testCase.ReqBody,
				Tags:  testCase.ReqTags,
			}))
			t.Is(testCase.ResStatusCode, res.StatusCode())
			t.Is(testCase.ResBody, res.Body())
		})
	}
}

func putArticles(articleID string, headers map[string]string, requestBody io.ReadCloser) *testing.Response {
	request := testing.PUT("/v1/articles/"+articleID, requestBody, &testing.RequestOptions{
		Headers: headers,
	})

	return testing.Do(request, router)
}
