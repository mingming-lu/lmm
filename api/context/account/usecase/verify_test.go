package usecase

import (
	"lmm/api/context/account/domain/service"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
)

func TestVerify_Success(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	origin := testingService.NewUser()
	token := service.EncodeToken(origin.Token)
	user, err := VerifyToken(token)

	tester.NoError(err)
	tester.Is(origin.ID, user.ID)
}

func TestVerify_InvalidToken(t *testing.T) {
	tester := testing.NewTester(t)

	user, err := VerifyToken("invalid")
	tester.Error(ErrInvalidToken, err)
	tester.Nil(user)
}
