package ui

// func TestPutV1UsersPassword(tt *testing.T) {
//
// 	tt.Run("Success", func(tt *testing.T) {
// 		c := context.Background()
// 		t := testing.NewTester(tt)
//
// 		username := "u" + uuidutil.NewUUID()[:8]
// 		email := username + "@example.com"
// 		password := uuidutil.NewUUID()
//
// 		if !t.NoError(ui.appService.RegisterNewUser(c, command.Register{
// 			UserName:     username,
// 			EmailAddress: email,
// 			Password:     password,
// 		})) {
// 			t.Fatal("failed to create user")
// 		}
//
// 		res := changePasswordByUserName(user.Name(), changePasswordRequestBody{
// 			OldPassword: password,
// 			NewPassword: "neWp@s$w0rD",
// 		})
//
// 		t.Is(http.StatusNoContent, res.StatusCode())
//
// 		var userTokenAfterPasswordChanged string
// 		if err := mysql.QueryRow(c,
// 			`select token from user where name = ?`,
// 			user.Name()).Scan(&userTokenAfterPasswordChanged); err != nil {
// 			t.Fatal(err)
// 		}
//
// 		t.Not(user.RawToken(), userTokenAfterPasswordChanged)
// 	})
//
// 	tt.Run("Status4xx", func(tt *testing.T) {
// 		type Case struct {
// 			UserName    string
// 			OldPassword string
// 			NewPassword string
// 			StatusCode  int
// 			ResBody     string
// 		}
//
// 		user := testutil.NewUser(mysql)
// 		cases := map[string]Case{
// 			"NoSuchUser": Case{
// 				UserName:    user.Name() + "qqq",
// 				OldPassword: user.RawPassword(),
// 				NewPassword: "MayBe@ValidPassword",
// 				StatusCode:  http.StatusNotFound,
// 				ResBody:     domain.ErrNoSuchUser.Error(),
// 			},
// 			"WrongPassword": Case{
// 				UserName:    user.Name(),
// 				OldPassword: user.RawPassword() + "aaa",
// 				NewPassword: "MayBe@ValidPassword",
// 				StatusCode:  http.StatusUnauthorized,
// 				ResBody:     domain.ErrUserPassword.Error(),
// 			},
// 			"UserNameMiss": Case{
// 				UserName:    testutil.NewUser(mysql).Name(),
// 				OldPassword: user.RawPassword(),
// 				NewPassword: "MayBe@ValidPassword",
// 				StatusCode:  http.StatusUnauthorized,
// 				ResBody:     domain.ErrUserPassword.Error(),
// 			},
// 			"EmptyOldPassword": Case{
// 				UserName:    user.Name(),
// 				NewPassword: "MayBe@ValidPassword",
// 				StatusCode:  http.StatusUnauthorized,
// 				ResBody:     domain.ErrUserPassword.Error(),
// 			},
// 			"EmptyNewPassword": Case{
// 				UserName:    user.Name(),
// 				OldPassword: user.RawPassword(),
// 				StatusCode:  http.StatusBadRequest,
// 				ResBody:     domain.ErrUserPasswordEmpty.Error(),
// 			},
// 			"NewPasswordTooShort": Case{
// 				UserName:    user.Name(),
// 				OldPassword: user.RawPassword(),
// 				NewPassword: "short",
// 				StatusCode:  http.StatusBadRequest,
// 				ResBody:     domain.ErrUserPasswordTooShort.Error(),
// 			},
// 			"NewPasswordTooWeak": Case{
// 				UserName:    user.Name(),
// 				OldPassword: user.RawPassword(),
// 				NewPassword: "123456789",
// 				StatusCode:  http.StatusBadRequest,
// 				ResBody:     domain.ErrUserPasswordTooWeak.Error(),
// 			},
// 			"NewPasswordTooLong": Case{
// 				UserName:    user.Name(),
// 				OldPassword: user.RawPassword(),
// 				NewPassword: strings.Repeat("a", 251),
// 				StatusCode:  http.StatusBadRequest,
// 				ResBody:     domain.ErrUserPasswordTooLong.Error(),
// 			},
// 		}
//
// 		for testname, testcase := range cases {
// 			tt.Run(testname, func(tt *testing.T) {
// 				t := testing.NewTester(tt)
// 				res := changePasswordByUserName(testcase.UserName, changePasswordRequestBody{
// 					OldPassword: testcase.OldPassword,
// 					NewPassword: testcase.NewPassword,
// 				})
//
// 				t.Is(testcase.StatusCode, res.StatusCode())
// 				t.Is(testcase.ResBody, res.Body())
// 			})
// 		}
// 	})
// }
//
// func changePasswordByUserName(username string, formData changePasswordRequestBody) *testing.Response {
// 	req := testing.PUT("/v1/users/"+username+"/password", &testing.RequestOptions{
// 		FormData: testing.StructToRequestBody(formData),
// 	})
//
// 	return testing.DoRequest(req, router)
// }
