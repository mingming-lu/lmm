package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	os.Setenv("GCP_PROJECT_ID", "appengine-id")

	t.Run("GAE", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS(""))
		router.POST("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		cases := map[string]struct {
			Origin string
			Allow  bool
		}{
			"NakedDomain": {
				Origin: "https://appengine-id-appspot.com",
				Allow:  true,
			},
			"Subdomain/WWW": {
				Origin: "https://www-dot-appengine-id-appspot.com",
				Allow:  true,
			},
			"Subdomain/Manager": {
				Origin: "https://manager-dot-appengine-id.appspot.com",
				Allow:  true,
			},
			"subdomain/Multi": {
				Origin: "https://one-dot-two-dot-appengine-id-appspot.com",
				Allow:  true,
			},
			"NotHTTPS": {
				Origin: "http://appengine-id-appspot.com",
				Allow:  false,
			},
			"NotAllowOrigin": {
				Origin: "https://example.com",
				Allow:  false,
			},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/ping", nil)
				req.Header.Set("Origin", testCase.Origin)
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				assert.Equal(t, testCase.Allow, res.Header().Get("Access-Control-Allow-Origin") == testCase.Origin)
				if testCase.Allow {
					assert.Equal(t, "pong", res.Body.String())
				} else {
					assert.Equal(t, "", res.Body.String())
				}
			})
		}
	})

	t.Run("CustomDomain", func(t *testing.T) {
		router := gin.New()
		router.Use(CORS("example.com"))
		router.POST("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		cases := map[string]struct {
			Origin string
			Allow  bool
		}{
			"NakedDomain": {
				Origin: "https://example.com",
				Allow:  true,
			},
			"Subdomain/WWW": {
				Origin: "https://www.example.com",
				Allow:  true,
			},
			"Subdomain/API": {
				Origin: "https://api.example.com",
				Allow:  true,
			},
			"Subdomain/Multi": {
				Origin: "https://api.dev.example.com",
				Allow:  true,
			},
			"NotHTTPS": {
				Origin: "http://example.com",
				Allow:  false,
			},
			"NotAllowOrigin": {
				Origin: "https://appengine-id-appspot.com",
				Allow:  false,
			},
		}

		for testName, testCase := range cases {
			t.Run(testName, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/ping", nil)
				req.Header.Set("Origin", testCase.Origin)
				res := httptest.NewRecorder()

				router.ServeHTTP(res, req)

				assert.Equal(t, testCase.Allow, res.Header().Get("Access-Control-Allow-Origin") == testCase.Origin)
				if testCase.Allow {
					assert.Equal(t, "pong", res.Body.String())
				} else {
					assert.Empty(t, res.Body.String())
				}
			})
		}

	})
}
