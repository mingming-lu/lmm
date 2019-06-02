package ui

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	authUtil "lmm/api/pkg/auth"
	httpUtil "lmm/api/pkg/http"
	"lmm/api/service/user/application"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/domain"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
)

var (
	bearerAuthPattern = regexp.MustCompile(`^Bearer +(.+)$`)
)

// UI handles user apis
type UI struct {
	appService *application.Service
}

// NewUI creates a new UI pointer
func NewUI(appService *application.Service) *UI {
	return &UI{appService: appService}
}

// SignUp handles POST /v1/users
func (ui *UI) SignUp(c *gin.Context) {
	reqBody := signUpRequestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		httpUtil.BadRequest(c)
		return
	}

	userID, err := ui.appService.RegisterNewUser(c, command.Register{
		UserName:     reqBody.Name,
		EmailAddress: reqBody.Email,
		Password:     reqBody.Password,
	})

	originalError := errors.Cause(err)
	switch originalError {
	case nil:
		c.Header("Location", fmt.Sprintf("users/%d", userID))
		httpUtil.Response(c, http.StatusCreated, "Success")

	case
		domain.ErrInvalidUserName,
		domain.ErrUserPasswordEmpty,
		domain.ErrInvalidPassword,
		domain.ErrUserPasswordTooShort,
		domain.ErrUserPasswordTooWeak,
		domain.ErrUserPasswordTooLong,
		domain.ErrInvalidEmail:
		c.String(http.StatusBadRequest, originalError.Error())

	case domain.ErrUserNameAlreadyUsed:
		c.String(http.StatusConflict, domain.ErrUserNameAlreadyUsed.Error())

	default:
		log.Criticalf(c, err.Error())
	}
}

// BasicAuth middleware
func (ui *UI) BasicAuth(next gin.HandlerFunc) gin.HandlerFunc {
	pattern := regexp.MustCompile(`^Basic +(.+)$`)

	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		matched := pattern.FindStringSubmatch(authHeader)
		if len(matched) != 2 {
			next(c)
			return
		}

		b, err := base64.URLEncoding.DecodeString(matched[1])
		if err != nil {
			next(c)
			return
		}

		basicauth := basicAuth{}
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(&basicauth); err != nil {
			next(c)
			return
		}

		auth, err := ui.appService.BasicAuth(c, command.Login{
			UserName: basicauth.UserName,
			Password: basicauth.Password,
		})
		if err != nil {
			next(c)
			return
		}

		ctxWithAuth := authUtil.NewContext(c.Request.Context(), auth)
		c.Request = c.Request.WithContext(ctxWithAuth)

		next(c)
	}
}

// BearerAuth is a middleware of bearer auth
func (ui *UI) BearerAuth(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		matched := bearerAuthPattern.FindStringSubmatch(authHeader)
		if len(matched) != 2 {
			next(c)
			return
		}

		auth, err := ui.appService.BearerAuth(c, matched[1])
		if err != nil {
			next(c)
			return
		}

		ctxWithAuth := authUtil.NewContext(c.Request.Context(), auth)
		c.Request = c.Request.WithContext(ctxWithAuth)

		next(c)
	}
}

// Token handles /v1/auth/token
func (ui *UI) Token(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Basic ") {
		ui.BasicAuth(func(c *gin.Context) {
			auth, ok := authUtil.FromContext(c)
			if !ok {
				httpUtil.Unauthorized(c)
				return
			}
			c.JSON(http.StatusOK, accessTokenView{
				AccessToken: auth.Token,
			})
		})(c)
		return

	} else if strings.HasPrefix(authHeader, "Bearer ") {
		matched := bearerAuthPattern.FindStringSubmatch(authHeader)
		if len(matched) != 2 {
			httpUtil.Unauthorized(c)
			return
		}

		token, err := ui.appService.RefreshAccessToken(c, matched[1])
		if err != nil {
			log.Warningf(c, err.Error())
			httpUtil.Unauthorized(c)
			return
		}
		c.JSON(http.StatusOK, accessTokenView{
			AccessToken: token.Hashed(),
		})
		return
	}

	httpUtil.Unauthorized(c)
}

// ChangeUserPassword handles PUT /v1/user/:user/password
func (ui *UI) ChangeUserPassword(c *gin.Context) {
	requestBody := changePasswordRequestBody{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Warningf(c, err.Error())
		httpUtil.BadRequest(c)
		return
	}

	err := ui.appService.UserChangePassword(c, command.ChangePassword{
		User:        c.Param("user"),
		OldPassword: requestBody.OldPassword,
		NewPassword: requestBody.NewPassword,
	})

	originalError := errors.Cause(err)
	switch originalError {
	case nil:
		httpUtil.Response(c, http.StatusOK, "Success")

	case domain.ErrUserPassword:
		c.String(http.StatusUnauthorized, domain.ErrUserPassword.Error())

	case
		domain.ErrUserPasswordEmpty,
		domain.ErrUserPasswordTooShort,
		domain.ErrUserPasswordTooWeak,
		domain.ErrUserPasswordTooLong,
		domain.ErrInvalidPassword:
		c.String(http.StatusBadRequest, originalError.Error())

	case domain.ErrNoSuchUser:
		c.String(http.StatusNotFound, domain.ErrNoSuchUser.Error())

	default:
		log.Criticalf(c, err.Error())
	}
}
