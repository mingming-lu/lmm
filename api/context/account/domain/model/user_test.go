package model

import (
	"lmm/api/testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tester := testing.NewTester(t)

	now := time.Now()

	user := NewUser(uint64(1234), "username", "password", "token", now)
	tester.Isa(&User{}, user)
	tester.Is(uint64(1234), user.ID())
	tester.Is("username", user.Name())
	tester.Is("password", user.Password())
	tester.Is("token", user.Token())
	tester.Is(now, user.CreatedAt())
}

func TestUpdatePassword(t *testing.T) {
	tester := testing.NewTester(t)

	now := time.Now()
	user := NewUser(uint64(1234), "username", "password", "token", now)

	tester.Error(user.VerifyPassword("newPassword"), "NewUser does nothing to encrypt password")

	user.UpdatePassword("newPassword")

	tester.Not("password", user.Password())
	tester.Not("newPassword", user.Password())
	tester.NoError(user.VerifyPassword("newPassword"))
}
