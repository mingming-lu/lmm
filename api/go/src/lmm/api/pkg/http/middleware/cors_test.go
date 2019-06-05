package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	t.Run("Origin", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS("https://example.com", "https://www.example.com", "https://manager-dot-example.com"))
		router.POST("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		cases := map[string]struct {
			Origin                          string
			ExpectAccessControllAllowOrigin string
		}{
			"Domain": {
				Origin:                          "https://example.com",
				ExpectAccessControllAllowOrigin: "https://example.com",
			},
			"Subdomain/WWW": {
				Origin:                          "https://www.example.com",
				ExpectAccessControllAllowOrigin: "https://www.example.com",
			},
			"Subdomain/Manager": {
				Origin:                          "https://manager-dot-example.com",
				ExpectAccessControllAllowOrigin: "https://manager-dot-example.com",
			},
			"NotHTTPS": {
				Origin:                          "http://example.com",
				ExpectAccessControllAllowOrigin: "",
			},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/ping", nil)
				req.Header.Set("Origin", testCase.Origin)
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				assert.Equal(t,
					testCase.ExpectAccessControllAllowOrigin,
					res.Header().Get("Access-Control-Allow-Origin"),
				)
			})
		}
	})
}
