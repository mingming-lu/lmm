package token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akinaru-lu/errors"
)

var (
	key        = []byte(os.Getenv("LMM_API_TOKEN_KEY"))
	Expire     = int64(86400) // default 1 days
	ErrExpired = errors.New("time expired")
)

func Encode(src string) ([]byte, error) {
	expire := time.Now().Unix() + Expire

	src = fmt.Sprintf("%v:%s", expire, src)
	b := []byte(src)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	encoded := make([]byte, aes.BlockSize+len(b))
	iv := encoded[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encoded[aes.BlockSize:], b)

	return encoded, nil
}

func Decode(src []byte) (string, error) {
	encoded, err := hex.DecodeString(fmt.Sprintf("%x", src))
	if err != nil {
		return "", errors.Wrap(err, "input string is not base64 encoded")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(encoded) < aes.BlockSize {
		return "", errors.New("the length of input string should be larger than or equal to 16")
	}

	iv := encoded[:aes.BlockSize]
	encoded = encoded[aes.BlockSize:]
	decoded := make([]byte, len(encoded))

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decoded, encoded)

	params := strings.Split(string(decoded), ":")
	if len(params) != 2 {
		return "", errors.New("access token format invalid")
	}
	seconds, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "invalid timestamp: "+params[0])
	}
	if time.Now().Unix() > seconds {
		return "", ErrExpired
	}
	return params[1], nil
}
