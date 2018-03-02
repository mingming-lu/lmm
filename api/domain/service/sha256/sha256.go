package sha256

import (
	"crypto/sha256"
	"fmt"
)

func CheckSum(data []byte) []byte {
	digest := sha256.Sum256(data)
	return digest[:]
}

func Hex(data []byte) string {
	return fmt.Sprintf("%x", CheckSum(data))
}
