package presentation

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
)

var (
	bearerAuthPattern = regexp.MustCompile(`^Bearer +(.+)$`)
)

type GinRouterProvider struct {
	appService *application.Service
}

func NewGinRouterProvider(appService *application.Service) *GinRouterProvider {
	return &GinRouterProvider{appService: appService}
}

func (p *GinRouterProvider) Provide(router *gin.Engine) {
	router.POST("/v1/users", p.SignUp)
	router.PUT("/v1/users/:user/password", p.ChangeUserPassword)

	router.POST("/v1/auth/token", p.Token)
}

// SignUp handles POST /v1/users
func (p *GinRouterProvider) SignUp(c *gin.Context) {
	reqBody := signUpRequestBody{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		httpUtil.BadRequest(c)
		return
	}

	userID, err := p.appService.RegisterNewUser(c, command.Register{
		UserName:     reqBody.Name,
		EmailAddress: reqBody.Email,
		Password:     reqBody.Password,
	})
	if err != nil {
		httpUtil.LogWarnf(c, err.Error())
	}

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
		httpUtil.LogCritf(c, err.Error())
	}
}

// BasicAuth middleware
func (p *GinRouterProvider) BasicAuth(next gin.HandlerFunc) gin.HandlerFunc {
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
			httpUtil.LogWarnf(c, "error on encoding base64: %s", err)
			next(c)
			return
		}

		basicauth := basicAuth{}
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(&basicauth); err != nil {
			httpUtil.LogWarnf(c, "error on decoding basic auth json: %s", err)
			next(c)
			return
		}

		auth, err := p.appService.BasicAuth(c, command.Login{
			UserName: basicauth.UserName,
			Password: basicauth.Password,
		})
		if err != nil {
			httpUtil.LogWarnf(c, "error on calling BasicAuth app service: %s", err)
			next(c)
			return
		}

		ctxWithAuth := authUtil.NewContext(c.Request.Context(), auth)
		c.Request = c.Request.WithContext(ctxWithAuth)

		next(c)
	}
}

// BearerAuth is a middleware of bearer auth
func (p *GinRouterProvider) BearerAuth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	matched := bearerAuthPattern.FindStringSubmatch(authHeader)
	if len(matched) != 2 {
		c.Next()
		return
	}

	auth, err := p.appService.BearerAuth(c, matched[1])
	if err != nil {
		httpUtil.LogWarnf(c, "error on calling BearerAuth app service: %s", err)
		c.Next()
		return
	}

	ctxWithAuth := authUtil.NewContext(c.Request.Context(), auth)
	c.Request = c.Request.WithContext(ctxWithAuth)

	c.Next()
}

// Token handles POST /v1/auth/token
func (p *GinRouterProvider) Token(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Basic ") {
		p.BasicAuth(func(c *gin.Context) {
			auth, ok := httpUtil.AuthFromGinContext(c)
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

		token, err := p.appService.RefreshAccessToken(c, matched[1])
		if err != nil {
			httpUtil.LogWarnf(c, err.Error())
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
func (p *GinRouterProvider) ChangeUserPassword(c *gin.Context) {
	requestBody := changePasswordRequestBody{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		httpUtil.LogWarnf(c, err.Error())
		httpUtil.BadRequest(c)
		return
	}

	err := p.appService.UserChangePassword(c, command.ChangePassword{
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
		httpUtil.LogCritf(c, err.Error())
	}
}
