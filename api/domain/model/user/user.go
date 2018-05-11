package user

import (
	"lmm/api/domain/model"
	"lmm/api/domain/service/sha256"
	"lmm/api/domain/service/uuid"
	"time"
)

type Minimal struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// User model is a Entry
type User struct {
	model.Entry
	Name      string
	Password  string
	GUID      string
	Token     string
	CreatedAt time.Time
}

// New new a User model which is going to be saved to repository (User.ID is 0!!)
// If want to build a User model from existing data, just do this -> user := &User{ID: xxx, ...}
func New(name, password string) *User {
	token := uuid.New()
	guid := uuid.New()
	encodedPassword := sha256.Hex([]byte(guid + password)) // digest

	return &User{
		Name:      name,
		Password:  encodedPassword,
		GUID:      guid,
		Token:     token,
		CreatedAt: time.Now(),
	}
}
