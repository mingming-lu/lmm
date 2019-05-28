package ui

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"lmm/api/http"
	authUtil "lmm/api/pkg/auth"
	"lmm/api/service/user/application"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
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
	router.PUT("/v1/users/:user/password", ui.ChangeUserPassword)

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

func TestPutV1UsersPassword(t *testing.T) {
	username := "U" + uuidutil.NewUUID()[:8]
	password := uuidutil.NewUUID() + uuidutil.NewUUID()
	email := username + "@lmm.local"

	res := postV1Users(signUpRequestBody{
		Name:     username,
		Password: password,
		Email:    email,
	})

	if !assert.Equal(t, http.StatusCreated, res.Code) {
		t.Fatal("failed to create new user: " + res.Body.String())
	}

	t.Run("Success", func(t *testing.T) {
		res := putV1UsersPassword(username, changePasswordRequestBody{
			OldPassword: password,
			NewPassword: uuidutil.NewUUID() + uuidutil.NewUUID(),
		})

		assert.Equal(t, 204, res.Code)
	})

	t.Run("Failure", func(t *testing.T) {
		type Case struct {
			UserName    string
			OldPassword string
			NewPassword string
			StatusCode  int
			ResBody     string
		}

		cases := map[string]Case{
			"NoSuchUser": Case{
				UserName:    username + "a",
				OldPassword: password,
				NewPassword: "MayBe@ValidPassword",
				StatusCode:  http.StatusNotFound,
				ResBody:     domain.ErrNoSuchUser.Error(),
			},
			"WrongPassword": Case{
				UserName:    username,
				OldPassword: password + "aa",
				NewPassword: "MayBe@ValidPassword",
				StatusCode:  http.StatusUnauthorized,
				ResBody:     domain.ErrUserPassword.Error(),
			},
			"EmptyOldPassword": Case{
				UserName:    username,
				NewPassword: "MayBe@ValidPassword",
				StatusCode:  http.StatusUnauthorized,
				ResBody:     domain.ErrUserPassword.Error(),
			},
			"EmptyNewPassword": Case{
				UserName:    username,
				OldPassword: password,
				StatusCode:  http.StatusBadRequest,
				ResBody:     domain.ErrUserPasswordEmpty.Error(),
			},
			"NewPasswordTooShort": Case{
				UserName:    username,
				OldPassword: password,
				NewPassword: "short",
				StatusCode:  http.StatusBadRequest,
				ResBody:     domain.ErrUserPasswordTooShort.Error(),
			},
			"NewPasswordTooWeak": Case{
				UserName:    username,
				OldPassword: password,
				NewPassword: "123456789",
				StatusCode:  http.StatusBadRequest,
				ResBody:     domain.ErrUserPasswordTooWeak.Error(),
			},
			"NewPasswordTooLong": Case{
				UserName:    username,
				OldPassword: password,
				NewPassword: strings.Repeat("a", 251),
				StatusCode:  http.StatusBadRequest,
				ResBody:     domain.ErrUserPasswordTooLong.Error(),
			},
		}

		for testname, testcase := range cases {
			t.Run(testname, func(t *testing.T) {
				res := putV1UsersPassword(testcase.UserName, changePasswordRequestBody{
					OldPassword: testcase.OldPassword,
					NewPassword: testcase.NewPassword,
				})

				assert.Equal(t, testcase.StatusCode, res.Code)
				assert.Equal(t, testcase.ResBody, res.Body.String())
			})
		}
	})
}

func TestBasicAuth(t *testing.T) {
	username := "U" + uuidutil.NewUUID()[:8]
	password := uuidutil.NewUUID() + uuidutil.NewUUID()
	email := username + "@lmm.local"

	postUserRes := postV1Users(signUpRequestBody{
		Name:     username,
		Password: password,
		Email:    email,
	})

	if !assert.Equal(t, http.StatusCreated, postUserRes.Code) {
		t.Fatal("failed to create user: ", postUserRes.Body.String())
	}

	location := postUserRes.Header().Get("Location")
	matched := regexp.MustCompile(`users/(\d+)`).FindStringSubmatch(location)
	userIDStr := matched[1]

	router := http.NewRouter()
	router.GET("/", ui.BasicAuth(func(c http.Context) {
		auth, ok := authUtil.FromContext(c)
		if !ok {
			http.Unauthorized(c)
			return
		}
		t.Run("FromContext", func(t *testing.T) {
			assert.Equal(t, userIDStr, strconv.FormatInt(auth.ID, 10))
			assert.Equal(t, username, auth.Name)
			assert.Equal(t, model.Ordinary.Name(), auth.Role)
			assert.NotEmpty(t, auth.Token)
		})
		c.String(http.StatusOK, "OK")
	}))

	t.Run("Authorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		buf := new(bytes.Buffer)
		if !assert.NoError(t, json.NewEncoder(buf).Encode(basicAuth{
			UserName: username,
			Password: password,
		})) {
			t.Fatal("unexpected failure of json encoding")
		}
		req.Header.Set("Authorization", "Basic "+base64.URLEncoding.EncodeToString(buf.Bytes()))

		res := httptest.NewRecorder()

		router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "OK", res.Body.String())
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
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

func putV1UsersPassword(username string, body changePasswordRequestBody) *httptest.ResponseRecorder {
	b, err := json.Marshal(body)
	if err != nil {
		panic(errors.Wrap(err, "failed to decode to json"))
	}

	req := httptest.NewRequest("PUT", "/v1/users/"+username+"/password", bytes.NewReader(b))
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	return res
}
