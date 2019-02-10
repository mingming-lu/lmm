package service

import (
	"lmm/api/service/user/domain/model"
	"lmm/api/testing"

	"github.com/google/uuid"
)

func TestBcryptService(tt *testing.T) {
	encrypter := BcryptService{}

	tt.Run("Encrypt", func(tt *testing.T) {
		t := testing.NewTester(tt)

		rawText := uuid.New().String()
		pw, err := model.NewPassword(rawText)
		if err != nil {
			tt.Fatal(err)
		}

		hashed, err := encrypter.Encrypt(pw)
		t.NoError(err)

		tt.Run("Verify", func(tt *testing.T) {
			t := testing.NewTester(tt)

			t.True(encrypter.Verify(rawText, hashed))
			t.False(encrypter.Verify("wrong password", hashed))
		})
	})
}
