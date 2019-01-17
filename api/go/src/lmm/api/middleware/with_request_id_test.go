package middleware

import (
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/util/uuidutil"
)

func TestWithRequestID(tt *testing.T) {
	t := testing.NewTester(tt)

	uuid := uuidutil.New()
	sig := uuidutil.New()

	router := http.NewRouter()
	router.Use(WithRequestID)
	router.GET("/", func(c http.Context) {
		requestID, ok := c.Value(http.RequestIDContextKey).(string)
		t.True(ok)
		t.Is(uuid, requestID)
		c.String(http.StatusOK, sig)
	})

	t.Run("WithRequestID", func(_ *testing.T) {
		req := testing.GET("/", nil)
		req.Header.Set("X-Request-ID", uuid)
		res := testing.DoRequest(req, router)

		t.Is(http.StatusOK, res.StatusCode())
		t.Is(sig, res.Body())
	})

	t.Run("WithoutRequestID", func(_ *testing.T) {
		req := testing.GET("/", nil)
		res := testing.DoRequest(req, router)

		t.Is(http.StatusBadRequest, res.StatusCode())
		t.Is(http.StatusText(http.StatusBadRequest), res.Body())
	})
}
