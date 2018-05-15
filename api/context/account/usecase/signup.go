package usecase

import (
	"lmm/api/context/account/domain/model"
)

func (uc *Usecase) SignUp(name, password string) (uint64, error) {
	m := model.New(name, password)
	user, err := uc.repo.Save(m)
	if err != nil {
		key, _, ok := uc.repo.CheckErrorDuplicate(err.Error())
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
