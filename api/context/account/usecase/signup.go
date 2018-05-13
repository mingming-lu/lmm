package usecase

import (
	"encoding/json"
	"io"
	"lmm/api/context/account/domain/model"
)

func (uc *Usecase) SignUp(requestBody io.ReadCloser) (uint64, error) {
	auth := &Auth{}
	err := json.NewDecoder(requestBody).Decode(auth)
	if err != nil {
		return 0, ErrInvalidInput
	}
	m := model.New(auth.Name, auth.Password)
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
