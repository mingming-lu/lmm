package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/akinaru-lu/elesion"
	"github.com/stretchr/testify/assert"

	"lmm/api/db"
	"lmm/api/utils/httptest"
)

var router *elesion.Router
var targetUser *User

func setUp() {
	db.Init("lmm_test")

	router = elesion.New()
	router.GET("/users/:user", GetUser)
	router.POST("/users", NewUser)

	targetUser = NewTestUser()
}

func tearDown() {
	err := db.New().DropDatabase("lmm_test").Close()
	if err != nil {
		fmt.Println(err)
	}
}

func TestMain(m *testing.M) {
	var code int
	defer func() {
		os.Exit(code)
	}()
	setUp()
	defer tearDown()
	code = m.Run()
}

// This test case tests both New and Get
func TestNewTestUser(t *testing.T) {
	assert.Equal(t, "test", targetUser.Name)
	assert.Equal(t, "testy", targetUser.Nickname)
	assert.Equal(t, int64(1), targetUser.ID)
}

func TestNewUser_EmptyName(t *testing.T) {
	user := User{}
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	r := httptest.POST("/users", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.StatusCode(), w.Body())
}

func TestNewUser_DuplicateName(t *testing.T) {
	b, err := json.Marshal(targetUser)
	assert.NoError(t, err)

	r := httptest.POST("/users", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.StatusCode(), w.Body()) // TODO status code should be 409 Conflict

	_, err = getUser(2)
	assert.Error(t, err, err.Error())
}

func TestPost_ResponseLocation(t *testing.T) {
	user := User{
		Name:     "Van Darkholme",
		Nickname: "van sama",
	}

	b, err := json.Marshal(user)
	assert.NoError(t, err)

	r := httptest.POST("/users", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.StatusCode(), w.Body())
	assert.Regexp(t, `^\/users\/(\d)+$`, w.Header().Get("Location"))
}

func TestGetUser_InvalidID(t *testing.T) {
	r := httptest.GET("/users/a")
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.StatusCode())
	assert.Equal(t, "invalid id: a\n", w.Body())
}

func TestGetUser_NotExist(t *testing.T) {
	r := httptest.GET("/users/99999")
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.StatusCode(), w.Body())
}
