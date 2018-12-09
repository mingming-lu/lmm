package ui

import (
	"lmm/api/http"
	"lmm/api/service/user/application"
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
		http.Panic(c, err.Error())
	}
}
