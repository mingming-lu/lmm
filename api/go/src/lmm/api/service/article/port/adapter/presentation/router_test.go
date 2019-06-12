package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	jsonUtil "lmm/api/pkg/json"
	testUtil "lmm/api/pkg/testing"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/port/adapter/persistence"
	"lmm/api/util/stringutil"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	router    *gin.Engine
	dataStore *datastore.Client
)

func TestMain(m *testing.M) {
	c := context.Background()

	var err error
	dataStore, err = datastore.NewClient(c, "")
	if err != nil {
		panic("failed to connect to datastore: " + err.Error())
	}

	router = gin.New()

	repo := persistence.NewArticleDataStore(dataStore)
	ui := NewUI(repo, repo, repo)

	router.POST("/v1/articles", testUtil.BearerAuth(dataStore, ui.PostNewArticle))
	router.PUT("/v1/articles/:articleID", testUtil.BearerAuth(dataStore, ui.PutV1Articles))

	code := m.Run()

	dataStore.Close()

	os.Exit(code)
}

func TestPostV1Articles(t *testing.T) {
	c := context.Background()

	user := testUtil.NewUser(c, dataStore)

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
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusCreated,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"message": "Success"}),
			ResHeaders:    map[string]string{"Location": `^\/v1\/articles\/\w+$`},
		},
		"NoTags": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusCreated,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"message": "Success"}),
			ResHeaders:    map[string]string{"Location": `^\/v1\/articles\/\w+$`},
		},
		"Unauthorized": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    nil,
			ResStatusCode: http.StatusUnauthorized,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": http.StatusText(http.StatusUnauthorized)}),
			ResHeaders:    nil,
		},
		"TitleRequired": {
			ReqTitle:      nil,
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": errTitleRequired.Error()}),
			ResHeaders:    nil,
		},
		"BodyRequired": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       nil,
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": errBodyRequired.Error()}),
			ResHeaders:    nil,
		},
		"TagsRequired": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       nil,
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": errTagsRequired.Error()}),
			ResHeaders:    nil,
		},
		"EmptyTitle": {
			ReqTitle:      stringutil.Pointer(""),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrEmptyArticleTitle.Error()}),
			ResHeaders:    nil,
		},
		"InvalidTitle": {
			ReqTitle:      stringutil.Pointer("!@#$"),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrInvalidArticleTitle.Error()}),
			ResHeaders:    nil,
		},
		"LongTitle": {
			ReqTitle:      stringutil.Pointer(strings.Repeat("t", 141)),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrArticleTitleTooLong.Error()}),
			ResHeaders:    nil,
		},
	}

	for testName, testCase := range cases {
		t.Run(testName, func(t *testing.T) {
			res := postV1Articles(
				testCase.ReqHeaders,
				postArticleAdapter{
					Title: testCase.ReqTitle,
					Body:  testCase.ReqBody,
					Tags:  testCase.ReqTags,
				},
			)
			assert.Equal(t, testCase.ResStatusCode, res.Code, testName)
			assert.JSONEq(t, testCase.ResBody, res.Body.String(), testName)
			for k, v := range testCase.ResHeaders {
				assert.Regexp(t, v, res.Header().Get(k), testName)
			}
		})
	}
}

func TestPutArticles(t *testing.T) {
	c := context.Background()

	user := testUtil.NewUser(c, dataStore)

	res := postV1Articles(
		http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
		postArticleAdapter{
			Title: stringutil.Pointer("title"),
			Body:  stringutil.Pointer("body"),
			Tags:  []string{"tag"},
		},
	)

	if res.Code != http.StatusCreated {
		t.Fatal("failed to create test article data")
	}

	groups := regexp.MustCompile(`^/v1/articles/(.+)$`).FindStringSubmatch(res.Header().Get("Location"))
	articleID := groups[1]

	cases := map[string]struct {
		ArticleID     string
		ReqLinkName   string
		ReqTitle      *string
		ReqBody       *string
		ReqTags       []string
		ReqHeaders    http.Header
		ResStatusCode int
		ResBody       string
	}{
		"Success": {
			ArticleID:     articleID,
			ReqLinkName:   uuid.New().String(),
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       []string{"foo", "bar"},
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusOK,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"message": "Success"}),
		},
		"NoTags": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusOK,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"message": "Success"}),
		},
		"Unauthorized": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    nil,
			ResStatusCode: http.StatusUnauthorized,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": http.StatusText(http.StatusUnauthorized)}),
		},
		"NotAuthor": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + testUtil.NewUser(c, dataStore).AccessToken}},
			ResStatusCode: http.StatusNotFound,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrNoSuchArticle.Error()}),
		},
		"NotFound": {
			ArticleID:     "notfound",
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + testUtil.NewUser(c, dataStore).AccessToken}},
			ResStatusCode: http.StatusNotFound,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrNoSuchArticle.Error()}),
		},
		"InvalidArticleID": {
			ArticleID:     "!nv@lidchrcter$",
			ReqTitle:      stringutil.Pointer("dummy"),
			ReqBody:       stringutil.Pointer("dummy"),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + testUtil.NewUser(c, dataStore).AccessToken}},
			ResStatusCode: http.StatusNotFound,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrNoSuchArticle.Error()}),
		},
		"TitleRequired": {
			ArticleID:     articleID,
			ReqTitle:      nil,
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": errTitleRequired.Error()}),
		},
		"BodyRequired": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       nil,
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": errBodyRequired.Error()}),
		},
		"TagsRequired": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       nil,
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": errTagsRequired.Error()}),
		},
		"EmptyTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(""),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrEmptyArticleTitle.Error()}),
		},
		"InvalidTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer("!@#$"),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrInvalidArticleTitle.Error()}),
		},
		"LongTitle": {
			ArticleID:     articleID,
			ReqTitle:      stringutil.Pointer(strings.Repeat("t", 141)),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       jsonUtil.MustJSONify(jsonUtil.JSON{"error": domain.ErrArticleTitleTooLong.Error()}),
		},
	}

	for testName, testCase := range cases {
		t.Run(testName, func(t *testing.T) {
			res := putV1Articles(
				testCase.ArticleID,
				testCase.ReqHeaders,
				postArticleAdapter{
					LinkName: &testCase.ReqLinkName,
					Title:    testCase.ReqTitle,
					Body:     testCase.ReqBody,
					Tags:     testCase.ReqTags,
				},
			)
			assert.Equal(t, testCase.ResStatusCode, res.Code)
			assert.Equal(t, testCase.ResBody, res.Body.String())
		})
	}
}

func postV1Articles(header http.Header, body postArticleAdapter) *httptest.ResponseRecorder {
	b, err := json.Marshal(body)
	if err != nil {
		panic(errors.Wrap(err, "failed to decode to json"))
	}

	req := httptest.NewRequest("POST", "/v1/articles", bytes.NewReader(b))
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}

func putV1Articles(articleID string, header http.Header, body postArticleAdapter) *httptest.ResponseRecorder {
	b, err := json.Marshal(body)
	if err != nil {
		panic(errors.Wrap(err, "failed to decode to json"))
	}

	req := httptest.NewRequest("PUT", "/v1/articles/"+articleID, bytes.NewReader(b))
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}
