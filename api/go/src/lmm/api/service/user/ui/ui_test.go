package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"lmm/api/http"
	"lmm/api/service/user/application"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/infra/persistence"
	"lmm/api/service/user/infra/service"
	"lmm/api/util/uuidutil"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	router *http.Router
	ui     *UI
)

func TestMain(m *testing.M) {
	dataStore, err := datastore.NewClient(context.Background(), "")
	if err != nil {
		panic(err)
	}

	userRepo := persistence.NewUserDataStore(dataStore)
	userAppService := application.NewService(&service.BcryptService{}, userRepo, userRepo)
	ui = NewUI(userAppService)

	router = http.NewRouter()
	router.POST("/v1/users", ui.SignUp)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPostV1Users(t *testing.T) {
	username := "U" + uuidutil.NewUUID()[:8]
	password := uuidutil.NewUUID() + uuidutil.NewUUID()
	email := username + "@lmm.local"

	t.Run("Created", func(t *testing.T) {
		res := postV1Users(signUpRequestBody{
			Name:     username,
			Password: password,
			Email:    email,
		})

		assert.Equal(t, 201, res.Code)
		assert.Regexp(t, regexp.MustCompile(`users/\d+`), res.Header().Get("Location"))
	})

	t.Run("BadRequest", func(tt *testing.T) {
		generateNewName := func() string {
			return "U" + uuidutil.NewUUID()[:8]
		}

		cases := map[string]struct {
			UserName   string
			Email      string
			Password   string
			StatusCode int
			Body       string
		}{
			"InvalidUserName": {
				"1234", email, password, 400, domain.ErrInvalidUserName.Error(),
			},
			"DuplicateUserName": {
				username, email, password, 409, domain.ErrUserNameAlreadyUsed.Error(),
			},
			"EmptyEmail": {
				generateNewName(), "", password, 400, domain.ErrInvalidEmail.Error(),
			},
			"InvalidEmail": {
				generateNewName(), "example.com", password, 400, domain.ErrInvalidEmail.Error(),
			},
			"EmptyPassword": {
				generateNewName(), email, "", 400, domain.ErrUserPasswordEmpty.Error(),
			},
			"InvalidPassword": {
				generateNewName(), email, "不合法的密码", 400, domain.ErrInvalidPassword.Error(),
			},
			"ShortPassword": {
				generateNewName(), email, "qwert", 400, domain.ErrUserPasswordTooShort.Error(),
			},
			"LongPassword": {
				generateNewName(), email, strings.Repeat("s", 251), 400, domain.ErrUserPasswordTooLong.Error(),
			},
			"WeakPassword": {
				generateNewName(), email, "password", 400, domain.ErrUserPasswordTooWeak.Error(),
			},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(tt *testing.T) {
				res := postV1Users(signUpRequestBody{
					Name:     testCase.UserName,
					Email:    testCase.Email,
					Password: testCase.Password,
				})

				assert.Equal(t, testCase.StatusCode, res.Code)
				assert.Equal(t, testCase.Body, res.Body.String())
			})
		}
	})
}

func postV1Users(body signUpRequestBody) *httptest.ResponseRecorder {
	b, err := json.Marshal(body)
	if err != nil {
		panic(errors.Wrap(err, "failed to decode to json"))
	}

	req := httptest.NewRequest("POST", "/v1/users", bytes.NewReader(b))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}
