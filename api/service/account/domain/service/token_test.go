package service

import (
	"encoding/base64"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestDecodeToken_Success(t *testing.T) {
	tester := testing.NewTester(t)

	rawToken := uuid.New()
	encodedToken := EncodeToken(rawToken)
	token, err := DecodeToken(encodedToken)

	tester.NoError(err)
	tester.Is(rawToken, token)
}

func TestDecodeToken_NotBase64Encoded(t *testing.T) {
	tester := testing.NewTester(t)

	token, err := DecodeToken("1526787124:c531db63-8b94-47c7-8a18-908b7b0ad0e8")

	tester.Error(err)
	tester.Is(ErrNotBase64Encoded, err)
	tester.Is("", token)
}

func TestDecodeToken_InvalidTokenLength(t *testing.T) {
	tester := testing.NewTester(t)

	token, err := DecodeToken(base64.StdEncoding.EncodeToString([]byte("invalid")))

	tester.Error(err)
	tester.Is(ErrInvalidTokenLength, err)
	tester.Is("", token)
}

func TestDecodeToken_InvalidTokenFormat(t *testing.T) {
	tester := testing.NewTester(t)

	token, err := DecodeToken(base64.StdEncoding.EncodeToString([]byte("1526787124::c531db63-8b94-47c7-8a18-908b7b0ad0e8")))

	tester.Error(err)
	tester.Is(ErrInvalidTokenFormat, err)
	tester.Is("", token)
}

func TestDecodeToken_InvalidTimestamp(t *testing.T) {
	tester := testing.NewTester(t)

	token, err := DecodeToken(base64.StdEncoding.EncodeToString([]byte("2001-01-01:c531db63-8b94-47c7-8a18-908b7b0ad0e8")))

	tester.Error(err)
	tester.Is(ErrInvalidTokenFormat, err)
	tester.Is("", token)
}

func TestDecodeToken_TokenExpired(t *testing.T) {
	tester := testing.NewTester(t)

	TokenExpire = int64(-10000)

	rawToken := uuid.New()
	encodedToken := EncodeToken(rawToken)
	token, err := DecodeToken(encodedToken)

	TokenExpire = int64(86400)

	tester.Error(err)
	tester.Is(ErrTokenExpired, err)
	tester.Is("", token)
}
