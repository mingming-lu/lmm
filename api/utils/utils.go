package utils

import (
	"github.com/google/uuid"
	"encoding/base64"
)

func NewUUID() string {
	return uuid.New().String()
}

func ToBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func FromBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
