package ui

import (
	"bytes"
	"io"
	"os"

	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/util/testutil"
	"mime/multipart"
)

func TestV1PostPhotos(tt *testing.T) {
	user := testutil.NewUser(mysql)

	type Case struct {
		// req
		RequestOptions *testing.RequestOptions

		// res
		StatusCode   int
		ResponseBody string
	}

	authorizedHeader := http.Header{"Authorization": []string{"Bearer " + user.AccessToken()}}

	cases := map[string]Case{
		"Unauthorized": {
			StatusCode:   http.StatusUnauthorized,
			ResponseBody: http.StatusText(http.StatusUnauthorized),
		},
		"BadRequest/NotMultipart": {
			RequestOptions: &testing.RequestOptions{
				Headers: authorizedHeader,
			},
			StatusCode:   http.StatusBadRequest,
			ResponseBody: http.StatusText(http.StatusBadRequest),
		},
		"Created": {
			RequestOptions: func() *testing.RequestOptions {
				formData := &bytes.Buffer{}
				f, err := os.Open("./.test/good.jpg")
				if err != nil {
					panic(err)
				}
				defer f.Close()

				mw := multipart.NewWriter(formData)
				fw, err := mw.CreateFormFile("photo", f.Name())
				if err != nil {
					panic(err)
				}
				defer mw.Close()

				if _, err := io.Copy(fw, f); err != nil {
					panic(err)
				}

				headers := http.Header{
					"Authorization": []string{"Bearer " + user.AccessToken()},
					"Content-Type":  []string{mw.FormDataContentType()},
				}

				return &testing.RequestOptions{
					Headers:  headers,
					FormData: formData,
				}
			}(),
			StatusCode:   http.StatusCreated,
			ResponseBody: "uploaded",
		},
	}

	for testname, testcase := range cases {
		tt.Run(testname, func(tt *testing.T) {
			t := testing.NewTester(tt)
			res := handler.postAssetsPhotos(testcase.RequestOptions)
			t.Is(testcase.StatusCode, res.StatusCode())
			t.Is(testcase.ResponseBody, res.Body())
		})
	}
}
