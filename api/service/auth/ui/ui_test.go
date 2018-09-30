package ui

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"lmm/api/http"
	"lmm/api/service/auth/infra/persistence"
	"lmm/api/storage/db"
	"lmm/api/testing"
	"lmm/api/util/testingutil"
)

var (
	dbSrcName  = "root:@tcp(lmm-mysql:3306)/"
	dbName     = os.Getenv("DATABASE_NAME")
	connParams = "parseTime=true"
)

var (
	dbEngine db.DB
	router   *http.Router
)

func TestMain(m *testing.M) {
	dbEngine = db.NewMySQL(fmt.Sprintf("%s%s?%s", dbSrcName, dbName, connParams))
	userRepo := persistence.NewUserStorage(dbEngine)
	ui := NewUI(userRepo)

	router = http.NewRouter()
	router.POST("/v1/auth/login", ui.Login)

	code := m.Run()
	os.Exit(code)
}
func TestLogin(tt *testing.T) {
	t := testing.NewTester(tt)

	type basicAuth struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	username := "U" + uuid.New().String()[:8]
	password := "password123!"
	_, err := testingutil.NewUserUser(dbEngine, username, password)
	if !t.NoError(err) {
		t.FailNow()
	}

	t.Run("Success", func(_ *testing.T) {
		auth := basicAuth{UserName: username, Password: password}
		b, err := json.Marshal(auth)
		if !t.NoError(err) {
			t.FailNow()
		}
		b64 := base64.URLEncoding.EncodeToString(b)

		headers := make(map[string]string)
		headers["Authorization"] = "Basic " + b64

		res := login(headers)
		t.Is(http.StatusOK, res.StatusCode())
	})
}

func login(headers map[string]string) *testing.Response {
	request := testing.POST("/v1/auth/login", nil)

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	res := testing.NewResponse()
	router.ServeHTTP(res, request)

	return res
}
