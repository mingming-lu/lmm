package ui

import (
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/util/testutil"
)

func TestV1PostPhotos(tt *testing.T) {
	t := testing.NewTester(tt)

	user := testutil.NewUser(mysql)

	type Case struct {
		// req
		RequestOptions *testing.RequestOptions

		// res
		StatusCode   int
		ResponseBody string
	}

	cases := map[string]Case{
		"Unauthorized": {
			StatusCode:   http.StatusUnauthorized,
			ResponseBody: http.StatusText(http.StatusUnauthorized),
		},
		"BadRequest/RequestBodyTooLarge": {
			RequestOptions: &testing.RequestOptions{
				Headers: http.Header{
					"Authorization": []string{"Bearer " + user.AccessToken()},
				},
			},
			StatusCode:   http.StatusBadRequest,
			ResponseBody: errRequestBodyTooLarge.Error(),
		},
	}

	for testname, testcase := range cases {
		t.Run(testname, func(_ *testing.T) {
			res := postAssetsPhotos(testcase.RequestOptions)
			t.Is(testcase.StatusCode, res.StatusCode())
			t.Is(testcase.ResponseBody, res.Body())
		})
	}
}

func postAssetsPhotos(opts *testing.RequestOptions) *testing.Response {
	req := testing.POST("/v1/assets/photos", opts)
	return handler(req)
}
