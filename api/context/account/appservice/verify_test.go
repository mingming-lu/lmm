package appservice

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/domain/service"
	"lmm/api/testing"
)

func TestVerify_Success(t *testing.T) {
	tester := testing.NewTester(t)
	app := New(repository.New(testing.DB()))

	name, password := randomUserNameAndPassword()
	app.SignUp(name, password)
	user, _ := app.SignIn(name, password)

	sameUser, err := app.VerifyToken(user.Token())

	tester.NoError(err)
	tester.Isa(&model.User{}, sameUser)
	tester.Is(user.ID(), sameUser.ID())

	rawToken, err := service.DecodeToken(user.Token())
	tester.NoError(err)
	tester.Is(rawToken, sameUser.Token())
}

func TestVerify_InvalidToken(t *testing.T) {
	tester := testing.NewTester(t)

	user, err := New(repository.New(testing.DB())).VerifyToken("invalid")
	tester.Error(ErrInvalidToken, err)
	tester.Nil(user)
}
