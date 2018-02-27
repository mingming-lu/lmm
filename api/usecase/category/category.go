package category

import (
	model "lmm/api/domain/model/category"
	repo "lmm/api/domain/repository/category"
)

func Register(userID int64, name string) (int64, error) {
	return repo.Add(userID, name)
}

func Update(userID, categoryID int64, name string) error {
	return repo.Update(userID, categoryID, name)
}

func FetchByUser(userID int64) ([]model.Minimal, error) {
	return repo.ByUser(userID)
}
