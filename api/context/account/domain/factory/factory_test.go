package factory

import (
	"lmm/api/testing"
	"lmm/api/utils/uuid"

	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	tester := testing.NewTester(t)
	name, password := uuid.New()[:32], uuid.New()

	user, err := NewUser(name, password)
	tester.NoError(err)
	tester.Is(name, user.Name())
	tester.NoError(bcrypt.CompareHashAndPassword([]byte(user.Password()), []byte(password)))
}
