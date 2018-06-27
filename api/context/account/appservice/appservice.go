package appservice

import (
	"encoding/json"
	"io"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/domain/service"
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
	json.NewDecoder(requestBody).Decode(&auth)

	user, err := app.userService.Register(auth.Name, auth.Password)
	if err != nil {
		return 0, err
	}
	return user.ID(), nil
}

// SignIn is a usecase which users sign in with a account
func (app *AppService) SignIn(requestBody io.ReadCloser) (*model.User, error) {
	auth := Auth{}
	json.NewDecoder(requestBody).Decode(&auth)

	user, err := app.userService.Login(auth.Name, auth.Password)

	// // user, err := app.repo.FindByName(name)
	// if err != nil {
	// 	if err.Error() == storage.ErrNoRows.Error() {
	// 		return nil, ErrInvalidUserNameOrPassword
	// 	}
	// 	return nil, err
	// }
	if err != nil {
		return nil, err
	}

	// if user.VerifyPassword(password) != nil {
	// 	return nil, ErrInvalidUserNameOrPassword
	// }

	return model.NewUser(
		user.ID(),
		user.Name(),
		user.Password(),
		service.EncodeToken(user.Token()),
		user.CreatedAt(),
	), nil
}

func (app *AppService) VerifyToken(encodedToken string) (user *model.User, err error) {
	// token, err := service.DecodeToken(encodedToken)
	// if err != nil {
	// 	return nil, ErrInvalidToken
	// }

	// user, err = app.repo.FindByToken(token)
	// if err != nil {
	// 	return nil, ErrInvalidToken
	// }
	// return user, nil
	return nil, nil
}

func (app *AppService) BearerAuth(auth string) (*model.User, error) {
	matched := PatternBearerAuthorization.FindStringSubmatch(auth)
	if len(matched) != 2 {
		return nil, service.ErrInvalidAuthorization
	}
	token := matched[1]

	return app.userService.GetUserByHashedToken(token)
}
