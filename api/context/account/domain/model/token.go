package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	key                   = []byte(os.Getenv("LMM_API_TOKEN_KEY"))
	TokenExpire           = int64(86400)
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTimestamp   = errors.New("invalid timestamp")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrInvalidTokenLength = errors.New("invalid token length")
	ErrTokenExpired       = errors.New("token expired")
	ErrNotBase64Encoded   = errors.New("not base64 encoded")
)

func encodeToken(rawToken string) string {
	expire := time.Now().Unix() + TokenExpire

	formattedToken := fmt.Sprintf("%v:%s", expire, rawToken)
	b := []byte(formattedToken)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	encodedToken := make([]byte, aes.BlockSize+len(b))
	iv := encodedToken[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encodedToken[aes.BlockSize:], b)

	return base64.StdEncoding.EncodeToString(encodedToken)
}

func decodeToken(encodedToken string) (string, error) {
	encodedTokenB, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return "", ErrNotBase64Encoded
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(encodedTokenB) < aes.BlockSize {
		return "", ErrInvalidTokenLength
	}

	iv := encodedTokenB[:aes.BlockSize]
	encodedTokenB = encodedTokenB[aes.BlockSize:]
	decodedTokenB := make([]byte, len(encodedToken))

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decodedTokenB, encodedTokenB)

	params := strings.Split(string(decodedTokenB), ":")
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
