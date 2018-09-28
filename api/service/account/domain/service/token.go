package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	key                   = []byte(os.Getenv("LMM_API_TOKEN_KEY"))
	TokenExpire           = int64(86400)
	ErrInvalidTimestamp   = errors.New("Invalid timestamp")
	ErrInvalidTokenFormat = errors.New("Invalid token format")
	ErrInvalidTokenLength = errors.New("Invalid token length")
	ErrTokenExpired       = errors.New("Token expired")
	ErrNotBase64Encoded   = errors.New("Not base64 encoded")
)

func init() {
	if len(key) == 0 {
		log.Fatalln("token key not set")
	}
}

// EncodeToken convert a token string into base64({timestamp}:{token}) format
func EncodeToken(targetToken string) string {
	expire := time.Now().Unix() + TokenExpire

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

// DecodeToken parse a token string from base64({timestamp}:{token}) format to raw token
// panic if failed to parse base64 encoded string
func DecodeToken(targetToken string) (string, error) {
	encodedToken, err := base64.StdEncoding.DecodeString(targetToken)
	if err != nil {
		return "", ErrNotBase64Encoded
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(encodedToken) < aes.BlockSize {
		return "", ErrInvalidTokenLength
	}

	iv := encodedToken[:aes.BlockSize]
	encodedToken = encodedToken[aes.BlockSize:]
	decodedToken := make([]byte, len(encodedToken))

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decodedToken, encodedToken)

	params := strings.Split(string(decodedToken), ":")
	if len(params) != 2 {
		return "", ErrInvalidTokenFormat
	}

	seconds, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return "", ErrInvalidTimestamp
	}

	if time.Now().Unix() > seconds {
		return "", ErrTokenExpired
	}
	return params[1], nil
}
