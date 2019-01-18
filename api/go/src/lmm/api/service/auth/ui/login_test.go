package ui

import (
	"encoding/base64"
	"encoding/json"

	"lmm/api/http"
	"lmm/api/service/auth/domain"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

func TestLogin(tt *testing.T) {
	user := testutil.NewUser(dbEngine)

	cases := map[string]struct {
		AuthFunc  func() string
		GrantType string
	}{
		"BasicAuth": {
			AuthFunc: func() string {
				type basicAuth struct {
					UserName string `json:"username"`
					Password string `json:"password"`
				}
				auth := basicAuth{UserName: user.Name(), Password: user.RawPassword()}
				b, err := json.Marshal(auth)
				if err != nil {
					tt.Fatal(err)
				}
				b64 := base64.URLEncoding.EncodeToString(b)

				return "Basic " + b64
			},
			GrantType: domain.GrantTypeBasicAuth,
		},
		"RefreshToken": {
			AuthFunc: func() string {
				return "Bearer " + user.AccessToken()
			},
			GrantType: domain.GrantTypeRefreshToken,
		},
	}

	for testname, testcase := range cases {
		tt.Run(testname, func(tt *testing.T) {
			t := testing.NewTester(tt)
			res := login(testcase.GrantType,
				http.Header{"Authorization": []string{testcase.AuthFunc()}})
			t.Is(http.StatusOK, res.StatusCode())

			accessToken := testutil.ExtractAccessToken(res.Body())

			tt.Run(testname+"/ValidateAccessToken", func(tt *testing.T) {
				t := testing.NewTester(tt)
				res := dummy(http.Header{"Authorization": []string{"Bearer " + accessToken}})
				t.Is(http.StatusOK, res.StatusCode())
				t.Is("OK", res.Body())
			})
		})
	}
}

func login(grantType string, headers http.Header) *testing.Response {
	request := testing.POST("/v1/auth/login", &testing.RequestOptions{
		Headers: headers,
		FormData: testing.StructToRequestBody(loginRequestBody{
			GrantType: grantType,
		}),
	})
	return testing.DoRequest(request, router)
}

func dummy(headers http.Header) *testing.Response {
	request := testing.GET("/v1/dummy", &testing.RequestOptions{
		Headers: headers,
	})

	return testing.DoRequest(request, router)
}
