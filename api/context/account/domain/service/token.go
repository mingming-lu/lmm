package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

// DecodeToken parse a token string from base64({token}:{timestamp}) format to raw token
// panic if failed to parse base64 encoded string
func DecodeToken(targetToken string) (string, error) {
	encodedToken, err := base64.StdEncoding.DecodeString(targetToken)
	if err != nil {
		panic("Failed to parse base64 encoded token: " + err.Error())
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(encodedToken) < aes.BlockSize {
		panic(fmt.Sprintf("Token size error: %d", len(encodedToken)))
	}

	iv := encodedToken[:aes.BlockSize]
	encodedToken = encodedToken[aes.BlockSize:]
	decodedToken := make([]byte, len(encodedToken))

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decodedToken, encodedToken)

	params := strings.Split(string(decodedToken), ":")
	if len(params) != 2 {
		panic("Invald token format: " + string(decodedToken))
	}

	seconds, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		panic("Invalid timestamp format: " + params[0])
	}

	if time.Now().Unix() > seconds {
		return "", ErrExpired
	}
	return params[1], nil
}
