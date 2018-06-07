package appservice

import (
	"lmm/api/context/account/domain/factory"
	"lmm/api/domain/repository"
)

func (uc *Usecase) SignUp(name, password string) (uint64, error) {
	if name == "" || password == "" {
		return 0, ErrEmptyUserNameOrPassword
	}

	m := factory.NewUser(name, password)
	user, err := uc.repo.Put(m)
	if err != nil {
		key, _, ok := repository.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "name" {
			return 0, ErrDuplicateUserName
		}
		return 0, err
	}
	return user.ID, nil
}
