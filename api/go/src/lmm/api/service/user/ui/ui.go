package ui

import (
	"lmm/api/http"
	"lmm/api/service/user/application"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/domain"

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

// ViewAllUsers handles GET /v1/users
func (ui *UI) ViewAllUsers(c http.Context) {
}
