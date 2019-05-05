package ui

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"lmm/api/http"
	authApp "lmm/api/service/auth/application"
	authStore "lmm/api/service/auth/infra/persistence/datastore"
	authUI "lmm/api/service/auth/ui"
	"lmm/api/service/user/application"
	"lmm/api/service/user/domain"
	userStore "lmm/api/service/user/infra/persistence"
	"lmm/api/testing"
	transaction "lmm/api/transaction/datastore"
	"lmm/api/util/stringutil"

	"cloud.google.com/go/datastore"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var (
	ui     *UI
	router *http.Router
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if err != nil {
		log.Fatalf(`failed to setup datastore: "%s"`, err.Error())
	}

	txManager := transaction.NewTransactionManager(client)

	authRepo := authStore.NewUserStore(client)
	authAppService := authApp.NewService(authRepo)
	authUI := authUI.NewUI(authAppService)

	userRepo := userStore.NewUserStore(client)
	appService := application.NewService(txManager, userRepo)
	ui = NewUI(appService)
	router = http.NewRouter()

	router.POST("/v1/users", ui.SignUp)
	router.PUT("/v1/users/:user/role", authUI.BearerAuth(ui.AssignUserRole))
	router.PUT("/v1/users/:user/password", ui.ChangeUserPassword)
	router.GET("/v1/users", authUI.BearerAuth(ui.ViewAllUsers))

	code := m.Run()

	os.Exit(code)
}

func TestPostUser(tt *testing.T) {
	username := "U" + stringutil.ReplaceAll(uuid.New().String(), "-", "")[:8]
	email := username + "@lmm.local"
	password := uuid.New().String()

	tt.Run("Success", func(tt *testing.T) {
		t := testing.NewTester(tt)
		res := postUser(testing.StructToRequestBody(signUpRequestBody{
			Name:     username,
			Email:    email,
			Password: password,
		}))

		t.Is(201, res.StatusCode())
		t.Is("success", res.Body())
	})

	tt.Run("Fail", func(tt *testing.T) {
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
				username, "", password, 400, domain.ErrInvalidEmail.Error(),
			},
			"InvalidEmail": {
				username, "example.com", password, 400, domain.ErrInvalidEmail.Error(),
			},
			"EmptyPassword": {
				username, email, "", 400, domain.ErrUserPasswordEmpty.Error(),
			},
			"InvalidPassword": {
				username, email, "不合法的密码", 400, domain.ErrInvalidPassword.Error(),
			},
			"ShortPassword": {
				username, email, "qwert", 400, domain.ErrUserPasswordTooShort.Error(),
			},
			"LongPassword": {
				username, email, strings.Repeat("s", 251), 400, domain.ErrUserPasswordTooLong.Error(),
			},
			"WeakPassword": {
				username, email, "password", 400, domain.ErrUserPasswordTooWeak.Error(),
			},
		}

		for testName, testCase := range cases {
			tt.Run(testName, func(tt *testing.T) {
				t := testing.NewTester(tt)
				res := postUser(testing.StructToRequestBody(signUpRequestBody{
					Name:     testCase.UserName,
					Email:    testCase.Email,
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

func getUsers(opts *testing.RequestOptions) *testing.Response {
	req := testing.GET("/v1/users", opts)

	return testing.DoRequest(req, router)
}

func assignUserRole(username string, opts *testing.RequestOptions) *testing.Response {
	req := testing.PUT("/v1/users/"+username+"/role", opts)

	return testing.DoRequest(req, router)
}
