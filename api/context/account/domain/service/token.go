package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/akinaru-lu/errors"
)

var (
	key = []byte(os.Getenv("LMM_API_TOKEN_KEY"))
	// Expire defines the expiration date of token, 1 day by default
	Expire = int64(86400)
	// ErrExpired should be return if token is expired
	ErrExpired = errors.New("time expired")
)

// EncodeToken convert a token string into base64({token}:{timestamp}) format
func EncodeToken(targetToken string) string {
	expire := time.Now().Unix() + Expire

	targetToken = fmt.Sprintf("%v:%s", expire, targetToken)
	b := []byte(targetToken)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	encoded := make([]byte, aes.BlockSize+len(b))
	iv := encoded[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encoded[aes.BlockSize:], b)

	return base64.StdEncoding.EncodeToString(encoded)
}
