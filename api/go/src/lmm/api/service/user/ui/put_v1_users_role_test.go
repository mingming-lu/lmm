package ui

// func TestPUTV1UsersRole(tt *testing.T) {
// 	tt.Run("Status4xx", func(tt *testing.T) {
// 		admin := testutil.NewAdmin(mysql)
// 		user := testutil.NewUser(mysql)
//
// 		authorizationHeader := http.Header{}
// 		authorizationHeader.Add("Authorization", "Bearer "+admin.AccessToken())
//
// 		type failCase struct {
// 			TargetUserName string
// 			ReqHeader      http.Header
// 			ReqBody        io.Reader
// 			ResStatus      int
// 			ResBody        string
// 		}
//
// 		type validForm struct {
// 			Role string `json:"role"`
// 		}
//
// 		testcases := map[string]failCase{
// 			"Unauthorized/NoAuthorizationHeader": failCase{
// 				TargetUserName: user.Name(),
// 				ResStatus:      http.StatusUnauthorized,
// 				ResBody:        http.StatusText(http.StatusUnauthorized),
// 			},
// 			"Unauthorized/NotAnAdmin": failCase{
// 				ReqHeader: func() http.Header {
// 					anotherUser := testutil.NewUser(mysql)
// 					headers := http.Header{}
// 					headers.Add("Authorization", "Bearer "+anotherUser.AccessToken())
// 					return headers
// 				}(),
// 				TargetUserName: user.Name(),
// 				ReqBody:        testing.StructToRequestBody(validForm{Role: "admin"}),
// 				ResStatus:      http.StatusForbidden,
// 				ResBody:        domain.ErrNoPermission.Error(),
// 			},
// 			"BadRequest/InvalidFormData": failCase{
// 				ReqHeader:      authorizationHeader,
// 				TargetUserName: user.Name(),
// 				ResStatus:      http.StatusBadRequest,
// 				ResBody:        http.StatusText(http.StatusBadRequest),
// 			},
// 			"BadRequest/NoSuchRole": failCase{
// 				ReqHeader:      authorizationHeader,
// 				TargetUserName: user.Name(),
// 				ReqBody:        testing.StructToRequestBody(validForm{Role: "dummy"}),
// 				ResStatus:      http.StatusBadRequest,
// 				ResBody:        domain.ErrNoSuchRole.Error(),
// 			},
// 			"BadRequest/CannotSelfAssign": failCase{
// 				ReqHeader:      authorizationHeader,
// 				TargetUserName: admin.Name(),
// 				ReqBody:        testing.StructToRequestBody(validForm{Role: "admin"}),
// 				ResStatus:      http.StatusBadRequest,
// 				ResBody:        domain.ErrCannotAssignSelfRole.Error(),
// 			},
// 			"NotFound/NoSuchUser": failCase{
// 				ReqHeader:      authorizationHeader,
// 				TargetUserName: "u" + uuid.New().String()[:7],
// 				ReqBody:        testing.StructToRequestBody(validForm{Role: "admin"}),
// 				ResStatus:      http.StatusNotFound,
// 				ResBody:        domain.ErrNoSuchUser.Error(),
// 			},
// 		}
//
// 		for testname, testcase := range testcases {
// 			tt.Run(testname, func(tt *testing.T) {
// 				t := testing.NewTester(tt)
//
// 				res := assignUserRole(testcase.TargetUserName, &testing.RequestOptions{
// 					FormData: testcase.ReqBody,
// 					Headers:  testcase.ReqHeader,
// 				})
// 				t.Is(testcase.ResStatus, res.StatusCode())
// 				t.Is(testcase.ResBody, res.Body())
// 			})
// 		}
// 	})
// }
