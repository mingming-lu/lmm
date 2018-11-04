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
	t := testing.NewTester(tt)
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
				if !t.NoError(err) {
					t.FailNow()
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
		t.Run(testname, func(_ *testing.T) {
			res := login(testcase.GrantType,
				map[string]string{"Authorization": testcase.AuthFunc()})
			t.Is(http.StatusOK, res.StatusCode())

			accessToken := testutil.ExtractAccessToken(res.Body())

			t.Run(testname+"/ValidateAccessToken", func(_ *testing.T) {
				res := dummy(map[string]string{"Authorization": "Bearer " + accessToken})
				t.Is(http.StatusOK, res.StatusCode())
				t.Is("OK", res.Body())
			})
		})
	}
}

func login(grantType string, headers map[string]string) *testing.Response {
	request := testing.POST(
		"/v1/auth/login",
		testing.StructToRequestBody(loginRequestBody{
			GrantType: grantType,
		}),
		&testing.RequestOptions{
			Headers: headers,
		},
	)
	return testing.Do(request, router)
}

func dummy(headers map[string]string) *testing.Response {
	request := testing.GET("/v1/dummy", &testing.RequestOptions{
		Headers: headers,
	})

	return testing.Do(request, router)
}
