package ui

import (
	"encoding/base64"
	"encoding/json"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

func TestLogin(tt *testing.T) {
	t := testing.NewTester(tt)

	type basicAuth struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	user := testutil.NewUser(dbEngine)

	t.Run("Success", func(_ *testing.T) {
		auth := basicAuth{UserName: user.Name(), Password: user.RawPassword()}
		b, err := json.Marshal(auth)
		if !t.NoError(err) {
			t.FailNow()
		}
		b64 := base64.URLEncoding.EncodeToString(b)

		headers := make(map[string]string)
		headers["Authorization"] = "Basic " + b64

		res := login(headers)
		t.Is(http.StatusOK, res.StatusCode())

		t.Run("BearerAuthOK", func(_ *testing.T) {
			accessToken := testutil.ExtractAccessToken(res.Body())
			t.NoError(err)
			t.Not("", accessToken)

			headers := make(map[string]string)
			headers["Authorization"] = "Bearer " + accessToken

			res := dummy(headers)
			t.Is(http.StatusOK, res.StatusCode())
			t.Is("OK", res.Body())
		})
	})
}

func login(headers map[string]string) *testing.Response {
	request := testing.POST("/v1/auth/login", nil, &testing.RequestOptions{
		Headers: headers,
	})
	return testing.Do(request, router)
}

func dummy(headers map[string]string) *testing.Response {
	request := testing.GET("/v1/dummy", &testing.RequestOptions{
		Headers: headers,
	})

	return testing.Do(request, router)
}
