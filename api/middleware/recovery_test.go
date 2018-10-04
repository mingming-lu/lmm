package middleware

import (
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"lmm/api/http"
	"lmm/api/storage/db"
	"lmm/api/testing"
)

func TestHandleTimeout(tt *testing.T) {
	t := testing.NewTester(tt)

	mysql := db.DefaultMySQL()
	defer mysql.Close()

	logger := zap.NewNop()
	defer logger.Sync()

	router := http.NewRouter()
	router.Use(NewRecovery(logger))
	router.GET("/timeout", func(c http.Context) {
		_, err := mysql.Exec(c, `select benchmark(2147483647, SHA2('slow down', 256))`)
		if err != nil {
			panic(err)
		}
	})

	res := testing.Do(testing.GET("/timeout", nil), router)
	t.Is(http.StatusRequestTimeout, res.StatusCode())
}
