package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akinaru-lu/elesion"
	"github.com/stretchr/testify/assert"
	"lmm/api/db"
	"lmm/api/utils/httptest"
	"net/http"
	"strings"
	"testing"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
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
	Name: "test name",
	Nickname: "test nickname",
}

func TestNewUser(t *testing.T) {
	b, err := json.Marshal(targetUser)
	assert.NoError(t, err)

	router := elesion.New()
	router.POST("/user", NewUser)
	r := httptest.POST("/user", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusCreated, w.StatusCode())
}

func TestNewUser_EmptyName(t *testing.T) {
	user := User{}
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	router := elesion.New()
	router.POST("/user", NewUser)
	r := httptest.POST("/user", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.StatusCode())
}

func TestNewUser_DuplicateName(t *testing.T) {
	user := User{
		Name: "test name",
	}
	b, err := json.Marshal(user)
	assert.NoError(t, err)

	router := elesion.New()
	router.POST("/user", NewUser)
	r := httptest.POST("/user", bytes.NewReader(b))
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.StatusCode()) // TODO status code should be 409 Conflict
}

func TestGetUser(t *testing.T) {
	router := elesion.New()
	router.GET("/user/:id", GetUser)
	r := httptest.GET("/user/1")
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.StatusCode())

	user := User{}
	err := json.NewDecoder(strings.NewReader(w.Body())).Decode(&user)
	assert.NoError(t, err)

	assert.Equal(t, targetUser.Name, user.Name)
	assert.Equal(t, targetUser.Nickname, user.Nickname)

	_, err = uuid.Parse(user.UID)
	assert.NoError(t, err)

	_, err = uuid.Parse(user.Token)
	assert.NoError(t, err)
}

func TestGetUser_InvalidID(t *testing.T) {
	router := elesion.New()
	router.GET("/user/:id", GetUser)
	r := httptest.GET("/user/a")
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.StatusCode())
	assert.Equal(t, "invalid id: a\n", w.Body())
}

func TestGetUser_NotExist(t *testing.T) {
	router := elesion.New()
	router.GET("/user/:id", GetUser)
	r := httptest.GET("/user/99999")
	w := httptest.NewResponseWriter()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.StatusCode())
}
