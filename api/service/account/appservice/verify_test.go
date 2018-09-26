package appservice

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/service"
	"lmm/api/context/account/infra"
	"lmm/api/testing"
)

func TestVerify_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	auth := randomUserNameAndPassword()

	_, err := app.SignUp(testing.StructToRequestBody(auth))
	t.NoError(err)

	user, _ := app.SignIn(testing.StructToRequestBody(auth))

	sameUser, err := app.VerifyToken(user.Token())

	t.NoError(err)
	t.Isa(&model.User{}, sameUser)
	t.Is(user.ID(), sameUser.ID())

	rawToken, err := service.DecodeToken(user.Token())
	t.NoError(err)
	t.Is(rawToken, sameUser.Token())
}

func TestVerify_InvalidToken(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewUserStorage(testing.DB())
	app := New(repo)

	user, err := app.VerifyToken("invalid token ???")
	t.IsError(service.ErrInvalidToken, err)
	t.Nil(user)
}
