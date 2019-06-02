package testing

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	c := context.Background()

	dataStore, err := datastore.NewClient(c, "")
	if err != nil {
		t.Fatal("failed to connect datastore: ", err.Error())
	}

	assert.NotPanics(t, func() {
		user := NewUser(c, dataStore)
		t.Logf("user created: %#v", user)

		assert.True(t, PasswordService.Verify(user.RawPassword, user.HashedPassword))

		token, err := TokenService.Decrypt(user.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, user.RawToken, token.Raw())
	})
}
