package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"lmm/api/http"
	testUtil "lmm/api/pkg/testing"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/infra/persistence"
	"lmm/api/util/stringutil"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	router    *http.Router
	dataStore *datastore.Client
)

func TestMain(m *testing.M) {
	c := context.Background()

	var err error
	dataStore, err = datastore.NewClient(c, "")
	if err != nil {
		panic("failed to connect to datastore: " + err.Error())
	}

	router = http.NewRouter()

	// TODO
	repo := persistence.NewArticleDataStore(dataStore)
	ui := NewUI(repo, repo, repo)

	router.POST("/v1/articles", testUtil.BearerAuth(dataStore, ui.PostNewArticle))

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
			ResBody:       "Success",
			ResHeaders:    map[string]string{"Location": `^\/v1\/articles\/\w+$`},
		},
		"NoTags": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
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
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTitleRequired.Error(),
			ResHeaders:    nil,
		},
		"BodyRequired": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       nil,
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errBodyRequired.Error(),
			ResHeaders:    nil,
		},
		"TagsRequired": {
			ReqTitle:      stringutil.Pointer(uuid.New().String()[:8]),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       nil,
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       errTagsRequired.Error(),
			ResHeaders:    nil,
		},
		"EmptyTitle": {
			ReqTitle:      stringutil.Pointer(""),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrEmptyArticleTitle.Error(),
			ResHeaders:    nil,
		},
		"InvalidTitle": {
			ReqTitle:      stringutil.Pointer("!@#$"),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrInvalidArticleTitle.Error(),
			ResHeaders:    nil,
		},
		"LongTitle": {
			ReqTitle:      stringutil.Pointer(strings.Repeat("t", 141)),
			ReqBody:       stringutil.Pointer(uuid.New().String()),
			ReqTags:       make([]string, 0),
			ReqHeaders:    http.Header{"Authorization": []string{"Bearer " + user.AccessToken}},
			ResStatusCode: http.StatusBadRequest,
			ResBody:       domain.ErrArticleTitleTooLong.Error(),
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
			assert.Equal(t, testCase.ResStatusCode, res.Code)
			assert.Equal(t, testCase.ResBody, res.Body.String())
			for k, v := range testCase.ResHeaders {
				assert.Regexp(t, v, res.Header().Get(k))
			}
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
