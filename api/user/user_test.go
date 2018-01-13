package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akinaru-lu/elesion"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"lmm/api/db"
	"lmm/api/utils/httptest"
	"net/http"
	"strings"
	"testing"
)

var router *elesion.Router

func TestMain(m *testing.M) {
	router = elesion.New()
	router.GET("/users/:user", GetUser)
	router.POST("/users", NewUser)

	db.Init("lmm_test")
	defer func() {
		err := db.New().DropDatabase("lmm_test").Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	m.Run()
}

var targetUser = User{
	Name:     "test name",
	Nickname: "test nickname",
}

func TestNewUser(t *testing.T) {
	b, err := json.Marshal(targetUser)
	assert.NoError(t, err)

	r := httptest.POST("/users", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.StatusCode(), w.Body())
	assert.Equal(t, w.Header().Get("Location"), "/users/1")
}

func TestNewUser_EmptyName(t *testing.T) {
	user := User{}
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	r := httptest.POST("/users", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.StatusCode())
}

func TestNewUser_DuplicateName(t *testing.T) {
	b, err := json.Marshal(targetUser)
	assert.NoError(t, err)

	r := httptest.POST("/users", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.StatusCode()) // TODO status code should be 409 Conflict
}

func TestGetUser(t *testing.T) {
	r := httptest.GET("/users/1")
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.StatusCode())

	user := User{}
	err := json.NewDecoder(strings.NewReader(w.Body())).Decode(&user)
	assert.NoError(t, err)

	assert.Equal(t, targetUser.Name, user.Name)
	assert.Equal(t, targetUser.Nickname, user.Nickname)

	_, err = uuid.Parse(user.GUID)
	assert.NoError(t, err)

	_, err = uuid.Parse(user.Token)
	assert.NoError(t, err)
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

	assert.Equal(t, http.StatusNotFound, w.StatusCode())
}
