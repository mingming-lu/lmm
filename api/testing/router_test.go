package testing

import (
	"time"

	"lmm/api/http"
)

func TestHandleTimeout(tt *T) {
	t := NewTester(tt)

	router := NewRouter()
	router.GET("/timeout", func(c http.Context) {
		time.Sleep(5 * time.Second)
	})

	res := Do(GET("/timeout", nil), router)
	t.Is(http.StatusRequestTimeout, res.StatusCode())
}
