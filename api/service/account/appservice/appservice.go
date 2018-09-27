package appservice

import (
	"encoding/json"
	"io"
	"lmm/api/service/account/domain/model"
	"lmm/api/service/account/domain/repository"
	"lmm/api/service/account/domain/service"
	"regexp"
)

var (
	PatternBearerAuthorization = regexp.MustCompile(`^Bearer (.+)$`)
)

type AppService struct {
	userService *service.UserService
}

type Auth struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func New(userRepo repository.UserRepository) *AppService {
	return &AppService{
		userService: service.NewUserService(userRepo),
	}
}

func (app *AppService) SignUp(requestBody io.ReadCloser) (uint64, error) {
	auth := Auth{}
	if json.NewDecoder(requestBody).Decode(&auth) != nil {
		return 0, service.ErrInvalidBody
	}

	user, err := app.userService.Register(auth.Name, auth.Password)
	if err != nil {
		return 0, err
	}
	return user.ID(), nil
}

// SignIn is a usecase which users sign in with a account
func (app *AppService) SignIn(requestBody io.ReadCloser) (*model.User, error) {
	auth := Auth{}
	if json.NewDecoder(requestBody).Decode(&auth) != nil {
		return nil, service.ErrInvalidBody
	}

	user, err := app.userService.Login(auth.Name, auth.Password)
	if err != nil {
		return nil, err
	}

	return model.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		service.EncodeToken(user.Token()),
		user.CreatedAt(),
	), nil
}

func (app *AppService) VerifyToken(hashedToken string) (user *model.User, err error) {
	user, err = app.userService.GetUserByHashedToken(hashedToken)
	if err != nil {
		return nil, service.ErrInvalidToken
	}
	return user, nil
}

func (app *AppService) BearerAuth(auth string) (*model.User, error) {
	matched := PatternBearerAuthorization.FindStringSubmatch(auth)
	if len(matched) != 2 {
		return nil, service.ErrInvalidAuthorization
	}
	token := matched[1]

	return app.userService.GetUserByHashedToken(token)
}