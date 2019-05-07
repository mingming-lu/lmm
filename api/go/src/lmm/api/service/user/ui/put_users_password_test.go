package ui

import (
	"net/http"
	"strings"

	"lmm/api/service/user/domain"
	"lmm/api/testing"
	"lmm/api/util/uuidutil"
)

func TestPutV1UsersPassword(tt *testing.T) {

	tt.Run("Success", func(tt *testing.T) {
		t := testing.NewTester(tt)

		username := "u" + uuidutil.NewUUID()[:8]
		email := username + "@example.com"
		password := uuidutil.NewUUID()

		if res := postUser(testing.StructToRequestBody(signUpRequestBody{
			Name:     username,
			Email:    email,
			Password: password,
		})); !t.Is(http.StatusCreated, res.StatusCode()) {
			t.Fatal("failed to create user")
		}

		res := changePasswordByUserName(username, changePasswordRequestBody{
			OldPassword: password,
			NewPassword: "neWp@s$w0rD",
		})

		t.Is(http.StatusNoContent, res.StatusCode())
	})

	tt.Run("Status4xx", func(tt *testing.T) {
		t := testing.NewTester(tt)

		type Case struct {
			UserName    string
			OldPassword string
			NewPassword string
			StatusCode  int
			ResBody     string
		}

		username := "u" + uuidutil.NewUUID()[:8]
		email := username + "@example.com"
		password := uuidutil.NewUUID()

		if res := postUser(testing.StructToRequestBody(signUpRequestBody{
			Name:     username,
			Email:    email,
			Password: password,
		})); !t.Is(http.StatusCreated, res.StatusCode()) {
			t.Fatal("failed to create user")
		}

		cases := map[string]Case{
			"NoSuchUser": Case{
				UserName:    username + "qqq",
				OldPassword: password,
				NewPassword: "MayBe@ValidPassword",
				StatusCode:  http.StatusNotFound,
				ResBody:     domain.ErrNoSuchUser.Error(),
			},
			"WrongPassword": Case{
				UserName:    username,
				OldPassword: password + "aaa",
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
			tt.Run(testname, func(tt *testing.T) {
				t := testing.NewTester(tt)
				res := changePasswordByUserName(testcase.UserName, changePasswordRequestBody{
					OldPassword: testcase.OldPassword,
					NewPassword: testcase.NewPassword,
				})

				t.Is(testcase.StatusCode, res.StatusCode())
				t.Is(testcase.ResBody, res.Body())
			})
		}
	})
}

func changePasswordByUserName(username string, formData changePasswordRequestBody) *testing.Response {
	req := testing.PUT("/v1/users/"+username+"/password", &testing.RequestOptions{
		FormData: testing.StructToRequestBody(formData),
	})

	return testing.DoRequest(req, router)
}
