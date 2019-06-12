package service

import (
	"testing"

	"lmm/api/service/user/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBcryptService(t *testing.T) {
	encrypter := BcryptService{}

	t.Run("Encrypt", func(t *testing.T) {
		rawText := uuid.New().String()
		pw, err := model.NewPassword(rawText)
		if err != nil {
			t.Fatal(err)
		}

		hashed, err := encrypter.Encrypt(pw)
		assert.NoError(t, err)

		t.Run("Verify", func(t *testing.T) {
			assert.True(t, encrypter.Verify(rawText, hashed))
			assert.False(t, encrypter.Verify("wrong password", hashed))
		})
	})
}
