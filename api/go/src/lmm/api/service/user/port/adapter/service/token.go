package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"
	"lmm/api/util/stringutil"

	"github.com/pkg/errors"
)

// NewTokenService returns a default token service implementation
func NewCFBTokenService(secretKey string, expire time.Duration) model.TokenService {
	return &cfbTokenService{
		secretKey: []byte(secretKey),
		expireSec: int64(expire.Seconds()),
	}
}

type cfbTokenService struct {
	secretKey []byte
	expireSec int64
}

func (s *cfbTokenService) Encrypt(rawToken string) (*model.AccessToken, error) {
	expire := time.Now().Unix() + s.expireSec

	token := fmt.Sprintf("%d:%s", expire, rawToken)
	b := []byte(token)

	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		panic(err.Error())
	}

	encoded := make([]byte, aes.BlockSize+len(b))
	iv := encoded[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err.Error())
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encoded[aes.BlockSize:], b)

	hashed := base64.URLEncoding.EncodeToString(encoded)

	return model.NewAccessToken(rawToken, hashed, time.Unix(expire, 0)), nil
}

func (s *cfbTokenService) Decrypt(hashed string) (*model.AccessToken, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(hashed)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		panic(err)
	}

	if len(tokenBytes) < aes.BlockSize {
		return nil, errors.Wrapf(domain.ErrInvalidTokenLength, "%d", len(tokenBytes))
	}

	iv, src := tokenBytes[:aes.BlockSize], tokenBytes[aes.BlockSize:]
	dst := make([]byte, len(src))

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(dst, src)

	params := strings.Split(string(dst), ":")
	if len(params) != 2 {
		return nil, domain.ErrInvalidTokenFormat
	}

	expire, err := stringutil.ParseInt64(params[0])
	if err != nil {
		return nil, errors.Wrap(domain.ErrInvalidTokenFormat, err.Error())
	}

	return model.NewAccessToken(params[1], hashed, time.Unix(expire, 0)), nil
}
