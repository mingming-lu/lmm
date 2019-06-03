package testing

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"

	"lmm/api/pkg/auth"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/infra/service"
	"lmm/api/util/uuidutil"

	"cloud.google.com/go/datastore"
)

var (
	// TokenService uses CFBTokenService as default
	TokenService = &service.CFBTokenService{}

	// PasswordService uses BscryptService as default
	PasswordService = &service.BcryptService{}
)

// User used for testing
type User struct {
	Key            *datastore.Key `datastore:"__key__"`
	Name           string         `datastore:"Name"`
	Role           string         `datastore:"Role"`
	RawPassword    string         `datastore:"Password"`
	HashedPassword string         `datastore:"-"`
	RawToken       string         `datastore:"Token"`
	AccessToken    string         `datastore:"-"`
	RegisteredAt   time.Time      `datastore:"RegisteredAt"`
}

// ID is a shortcut for user.Key.ID
func (user *User) ID() int64 {
	return user.Key.ID
}

// NewUser create a new user
func NewUser(ctx context.Context, dataStore *datastore.Client) *User {
	username := "U" + uuidutil.NewUUID()[:8]
	password := uuidutil.NewUUID() + uuidutil.NewUUID()

	hashedPassword, err := factory.NewFactory(PasswordService, nil).NewPassword(password)
	if err != nil {
		panic("failed to encrypt password: " + err.Error())
	}

	token := uuidutil.NewUUID()
	accessToken, err := TokenService.Encrypt(token)
	if err != nil {
		panic("failed to generate access token: " + err.Error())
	}

	user := &User{
		Name:           username,
		Role:           "Ordinary",
		RawPassword:    password,
		HashedPassword: hashedPassword,
		RawToken:       token,
		AccessToken:    accessToken.Hashed(),
		RegisteredAt:   time.Now(),
	}
	key, err := dataStore.Put(ctx, datastore.IncompleteKey("User", nil), user)
	if err != nil {
		panic("failed to put user: " + err.Error())
	}

	user.Key = key

	fmt.Printf("INFO: created user: %#v\n", user)

	return user
}

// BearerAuth middleware for testing
func BearerAuth(dataStore *datastore.Client, next gin.HandlerFunc) gin.HandlerFunc {
	pattern := regexp.MustCompile(`^Bearer +(.+)$`)

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		matched := pattern.FindStringSubmatch(authHeader)
		if len(matched) != 2 {
			log.Printf("invalid header: %s", authHeader)
			next(c)
			return
		}

		accessToken := matched[1]
		token, err := TokenService.Decrypt(accessToken)
		if err != nil {
			log.Print(err.Error())
			next(c)
			return
		}

		q := datastore.NewQuery("User").KeysOnly().Filter("Token =", token.Raw())
		keys, err := dataStore.GetAll(c, q, nil)
		if err != nil {
			log.Print(err.Error())
			next(c)
			return
		}

		users := make([]*User, len(keys))
		if err := dataStore.GetMulti(c, keys, users); err != nil {
			log.Print(err.Error())
			next(c)
			return
		}

		user := users[0]
		ctxWithAuth := auth.NewContext(c.Request.Context(), &auth.Auth{
			ID:    user.ID(),
			Name:  user.Name,
			Token: user.RawToken,
			Role:  user.Role,
		})
		c.Request = c.Request.WithContext(ctxWithAuth)

		next(c)
	}
}
