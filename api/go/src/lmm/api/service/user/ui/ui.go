package ui

import (
	"encoding/json"

	"lmm/api/http"
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

	_, err := ui.appService.RegisterNewUser(c,
		reqBody.Name,
		reqBody.Password,
	)

	switch errors.Cause(err) {
	case nil:
		c.String(http.StatusCreated, "success")
	case domain.ErrInvalidUserName:
		c.String(http.StatusBadRequest, domain.ErrInvalidUserName.Error())
	case domain.ErrUserNameAlreadyUsed:
		c.String(http.StatusConflict, domain.ErrUserNameAlreadyUsed.Error())
	case domain.ErrUserPasswordEmpty:
		c.String(http.StatusBadRequest, domain.ErrUserPasswordEmpty.Error())
	case domain.ErrInvalidPassword:
		c.String(http.StatusBadRequest, domain.ErrInvalidPassword.Error())
	case domain.ErrUserPasswordTooShort:
		c.String(http.StatusBadRequest, domain.ErrUserPasswordTooShort.Error())
	case domain.ErrUserPasswordTooLong:
		c.String(http.StatusBadRequest, domain.ErrUserPasswordTooLong.Error())
	case domain.ErrUserPasswordTooWeak:
		c.String(http.StatusBadRequest, domain.ErrUserPasswordTooWeak.Error())
	default:
		http.Log().Panic(c, err.Error())
	}
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
	userName := c.Request().Header.Get("X-LMM-ID")
	if userName == "" {
		http.Unauthorized(c)
		return
	}

	requestBody := changePasswordRequestBody{}
	if err := c.Request().Bind(&requestBody); err != nil {
		http.Log().Warn(c, err.Error())
		http.BadRequest(c)
		return
	}

	err := ui.appService.UserChangePassword(c, command.ChangePassword{
		User:        userName,
		OldPassword: requestBody.OldPassword,
		NewPassword: requestBody.NewPassword,
	})

	originalError := errors.Cause(err)
	switch originalError {
	case nil:
		http.NoContent(c)

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
