package service

import (
	"testing"

	"lmm/api/util/uuidutil"

	"github.com/stretchr/testify/assert"
)

func TestCFBTokenService(t *testing.T) {
	cfb := &CFBTokenService{}

	token := uuidutil.NewUUID()

	t.Run("Encrypt", func(t *testing.T) {
		accessToken, err := cfb.Encrypt(token)
		assert.NoError(t, err)
		assert.Equal(t, token, accessToken.Raw())
		assert.False(t, accessToken.Expired())

		t.Run("Decrypt", func(t *testing.T) {
			sameAccessToken, err := cfb.Decrypt(accessToken.Hashed())
			assert.NoError(t, err)
			assert.Equal(t, accessToken, sameAccessToken)
		})
	})
}
