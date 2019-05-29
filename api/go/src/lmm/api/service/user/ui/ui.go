package ui

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"lmm/api/http"
	authUtil "lmm/api/pkg/auth"
	"lmm/api/service/user/application"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/application/query"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/model"

	"github.com/pkg/errors"
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
func (ui *UI) SignUp(c http.Context) {
	reqBody := signUpRequestBody{}
	c.Request().Bind(&reqBody)

	userID, err := ui.appService.RegisterNewUser(c, command.Register{
		UserName:     reqBody.Name,
		EmailAddress: reqBody.Email,
		Password:     reqBody.Password,
	})

	originalError := errors.Cause(err)
	switch originalError {
	case nil:
		c.Header("Location", fmt.Sprintf("users/%d", userID))
		c.String(http.StatusCreated, "success")

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
		http.Log().Panic(c, err.Error())
	}
}

// BasicAuth middleware
func (ui *UI) BasicAuth(next http.Handler) http.Handler {
	pattern := regexp.MustCompile(`^Basic +(.+)$`)

	return func(c http.Context) {
		authHeader := c.Request().Header.Get("Authorization")

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

		ctx := authUtil.NewContext(c.Request().Context(), auth)

		next(c.With(ctx))
	}
}

// BearerAuth is a middleware of bearer auth
func (ui *UI) BearerAuth(next http.Handler) http.Handler {
	return func(c http.Context) {
		authHeader := c.Request().Header.Get("Authorization")

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

		ctx := authUtil.NewContext(c.Request().Context(), auth)

		next(c.With(ctx))
	}
}

var bearerAuthPattern = regexp.MustCompile(`^Bearer +(.+)$`)

// Token handles /v1/auth/token
func (ui *UI) Token(c http.Context) {
	authHeader := c.Request().Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Basic ") {
		ui.BasicAuth(func(c http.Context) {
			auth, ok := authUtil.FromContext(c)
			if !ok {
				http.Unauthorized(c)
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
			http.Unauthorized(c)
			return
		}

		token, err := ui.appService.RefreshAccessToken(c, matched[1])
		if err != nil {
			http.Log().Warn(c, err.Error())
			http.Unauthorized(c)
			return
		}
		c.JSON(http.StatusOK, accessTokenView{
			AccessToken: token.Hashed(),
		})
		return
	}

	http.Unauthorized(c)
}

// AssignUserRole handles PUT /v1/users/:user/role
func (ui *UI) AssignUserRole(c http.Context) {
	userName := c.Request().Header.Get("X-LMM-ID")
	if userName == "" {
		http.Unauthorized(c)
		return
	}

	targetUserName := c.Request().PathParam("user")
	requestBody := &assignRoleRequestBody{}
	if err := c.Request().Bind(&requestBody); err != nil {
		http.BadRequest(c)
		return
	}

	err := ui.appService.AssignRole(c, command.AssignRole{
		OperatorUser: userName,
		TargetUser:   targetUserName,
		TargetRole:   requestBody.Role,
	})
	switch errors.Cause(err) {
	case nil:
		http.NoContent(c)
	case domain.ErrCannotAssignSelfRole:
		c.String(http.StatusBadRequest, domain.ErrCannotAssignSelfRole.Error())
	case domain.ErrNoSuchRole:
		c.String(http.StatusBadRequest, domain.ErrNoSuchRole.Error())
	case domain.ErrNoPermission:
		c.String(http.StatusForbidden, domain.ErrNoPermission.Error())
	case domain.ErrNoSuchUser:
		c.String(http.StatusNotFound, domain.ErrNoSuchUser.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}

// ChangeUserPassword handles PUT /v1/user/:user/password
func (ui *UI) ChangeUserPassword(c http.Context) {
	requestBody := changePasswordRequestBody{}
	if err := c.Request().Bind(&requestBody); err != nil {
		http.Log().Warn(c, err.Error())
		http.BadRequest(c)
		return
	}

	err := ui.appService.UserChangePassword(c, command.ChangePassword{
		User:        c.Request().PathParam("user"),
		OldPassword: requestBody.OldPassword,
		NewPassword: requestBody.NewPassword,
	})

	originalError := errors.Cause(err)
	switch originalError {
	case nil:
		http.NoContent(c)

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
		http.Log().Panic(c, err.Error())
	}
}

// ViewAllUsers handles GET /v1/users
func (ui *UI) ViewAllUsers(c http.Context) {
	userName := c.Request().Header.Get("X-LMM-ID")
	if userName == "" {
		http.Unauthorized(c)
		return
	}

	q := query.ViewAllUsers{
		Page:    c.Request().QueryParamOrDefault("page", "1"),
		Count:   c.Request().QueryParamOrDefault("count", "50"),
		OrderBy: c.Request().QueryParamOrDefault("sort_by", "registered_date"),
		Order:   c.Request().QueryParamOrDefault("sort", "desc"),
	}
	users, totalUsers, err := ui.appService.ViewAllUsersByOptions(c, q)

	switch errors.Cause(err) {
	case nil:
		c.JSON(http.StatusOK, ui.usersToJSONView(q, users, totalUsers))
	case domain.ErrInvalidPage:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidCount:
		c.String(http.StatusBadRequest, err.Error())
	case domain.ErrInvalidViewOrder:
		c.String(http.StatusBadRequest, err.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
}

func (ui *UI) usersToJSONView(query query.ViewAllUsers, users []*model.UserDescriptor, totalUsers uint) usersView {
	userItems := make([]userView, len(users), len(users))
	for i, user := range users {
		userItems[i] = userView{
			Name:           user.Name(),
			Role:           user.Role().Name(),
			RegisteredDate: user.RegisteredAt().Unix(),
		}
	}

	return usersView{
		Users:  userItems,
		Count:  json.Number(query.Count),
		Page:   json.Number(query.Page),
		Total:  totalUsers,
		Sort:   query.Order,
		SortBy: query.OrderBy,
	}
}
