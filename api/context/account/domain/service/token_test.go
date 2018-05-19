package service

import (
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestDecodeToken_Success(t *testing.T) {
	// 1526787124:c531db63-8b94-47c7-8a18-908b7b0ad0e8
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
