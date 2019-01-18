package ui

import (
	"io"
	"os"
	"strings"

	"lmm/api/http"
	"lmm/api/service/user/application"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/infra/persistence"
	"lmm/api/storage/db"
	"lmm/api/testing"
	"lmm/api/util/stringutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var (
	router *http.Router
)

func TestMain(m *testing.M) {
	mysql := db.DefaultMySQL()
	userRepo := persistence.NewUserStorage(mysql)
	appService := application.NewService(userRepo)
	ui := NewUI(appService)
	router = http.NewRouter()
	router.POST("/v1/users", ui.SignUp)

	code := m.Run()

	if err := mysql.Close(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestPostUser(tt *testing.T) {
	username := "U" + stringutil.ReplaceAll(uuid.New().String(), "-", "")[:8]
	password := uuid.New().String()

	tt.Run("Success", func(tt *testing.T) {
		t := testing.NewTester(tt)
		res := postUser(testing.StructToRequestBody(signUpRequestBody{
			Name:     username,
			Password: password,
		}))

		t.Is(201, res.StatusCode())
		t.Is("success", res.Body())
	})

	tt.Run("Fail", func(tt *testing.T) {
		cases := map[string]struct {
			UserName   string
			Password   string
			StatusCode int
			Body       string
		}{
			"InvalidUserName": {
				"1234", password, 400, domain.ErrInvalidUserName.Error(),
			},
			"DuplicateUserName": {
				username, password, 409, domain.ErrUserNameAlreadyUsed.Error(),
			},
			"EmptyPassword": {
				username, "", 400, domain.ErrUserPasswordEmpty.Error(),
			},
			"InvalidPassword": {
				username, "不合法的密码", 400, domain.ErrInvalidPassword.Error(),
			},
			"ShortPassword": {
				username, "qwert", 400, domain.ErrUserPasswordTooShort.Error(),
			},
			"LongPassword": {
				username, strings.Repeat("s", 251), 400, domain.ErrUserPasswordTooLong.Error(),
			},
			"WeakPassword": {
				username, "password", 400, domain.ErrUserPasswordTooWeak.Error(),
			},
		}

		for testName, testCase := range cases {
			tt.Run(testName, func(tt *testing.T) {
				t := testing.NewTester(tt)
				res := postUser(testing.StructToRequestBody(signUpRequestBody{
					Name:     testCase.UserName,
					Password: testCase.Password,
				}))

				t.Is(testCase.StatusCode, res.StatusCode())
				t.Is(testCase.Body, res.Body())
			})
		}
	})
}

func postUser(requestBody io.ReadCloser) *testing.Response {
	request := testing.POST("/v1/users", &testing.RequestOptions{
		FormData: requestBody,
	})

	return testing.DoRequest(request, router)
}
