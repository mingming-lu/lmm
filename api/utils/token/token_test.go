package token

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"io"
	"crypto/rand"
	"strconv"
)


func TestMain(m *testing.M) {
	key = make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

var (
	accessToken string
	msgToEncode = "test message"
)

func TestEncode(t *testing.T) {
	token, err := Encode(msgToEncode)
	assert.NoError(t, err)
	accessToken = token
}

func TestDecode(t *testing.T) {
	time, msg, err := Decode(accessToken)
	assert.NoError(t, err)

	// if is timestamp in second
	_, err = strconv.ParseInt(time, 10, 64)
	assert.NoError(t, err)

	assert.Equal(t, msgToEncode, msg)
}

func TestRefresh(t *testing.T) {

}
