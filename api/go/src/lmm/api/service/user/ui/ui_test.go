package ui

// func TestPostUser(tt *testing.T) {
// 	username := "U" + stringutil.ReplaceAll(uuid.New().String(), "-", "")[:8]
// 	email := username + "@lmm.local"
// 	password := uuid.New().String()
//
// 	tt.Run("Success", func(tt *testing.T) {
// 		t := testing.NewTester(tt)
// 		res := postUser(testing.StructToRequestBody(signUpRequestBody{
// 			Name:     username,
// 			Email:    email,
// 			Password: password,
// 		}))
//
// 		t.Is(201, res.StatusCode())
// 		t.Is("success", res.Body())
// 	})
//
// 	tt.Run("Fail", func(tt *testing.T) {
// 		cases := map[string]struct {
// 			UserName   string
// 			Email      string
// 			Password   string
// 			StatusCode int
// 			Body       string
// 		}{
// 			"InvalidUserName": {
// 				"1234", email, password, 400, domain.ErrInvalidUserName.Error(),
// 			},
// 			"DuplicateUserName": {
// 				username, email, password, 409, domain.ErrUserNameAlreadyUsed.Error(),
// 			},
// 			"EmptyEmail": {
// 				username, "", password, 400, domain.ErrInvalidEmail.Error(),
// 			},
// 			"InvalidEmail": {
// 				username, "example.com", password, 400, domain.ErrInvalidEmail.Error(),
// 			},
// 			"EmptyPassword": {
// 				username, email, "", 400, domain.ErrUserPasswordEmpty.Error(),
// 			},
// 			"InvalidPassword": {
// 				username, email, "不合法的密码", 400, domain.ErrInvalidPassword.Error(),
// 			},
// 			"ShortPassword": {
// 				username, email, "qwert", 400, domain.ErrUserPasswordTooShort.Error(),
// 			},
// 			"LongPassword": {
// 				username, email, strings.Repeat("s", 251), 400, domain.ErrUserPasswordTooLong.Error(),
// 			},
// 			"WeakPassword": {
// 				username, email, "password", 400, domain.ErrUserPasswordTooWeak.Error(),
// 			},
// 		}
//
// 		for testName, testCase := range cases {
// 			tt.Run(testName, func(tt *testing.T) {
// 				t := testing.NewTester(tt)
// 				res := postUser(testing.StructToRequestBody(signUpRequestBody{
// 					Name:     testCase.UserName,
// 					Email:    testCase.Email,
// 					Password: testCase.Password,
// 				}))
//
// 				t.Is(testCase.StatusCode, res.StatusCode())
// 				t.Is(testCase.Body, res.Body())
// 			})
// 		}
// 	})
// }
//
// func postUser(requestBody io.ReadCloser) *testing.Response {
// 	request := testing.POST("/v1/users", &testing.RequestOptions{
// 		FormData: requestBody,
// 	})
//
// 	return testing.DoRequest(request, router)
// }
//
// func getUsers(opts *testing.RequestOptions) *testing.Response {
// 	req := testing.GET("/v1/users", opts)
//
// 	return testing.DoRequest(req, router)
// }
//
// func assignUserRole(username string, opts *testing.RequestOptions) *testing.Response {
// 	req := testing.PUT("/v1/users/"+username+"/role", opts)
//
// 	return testing.DoRequest(req, router)
// }
