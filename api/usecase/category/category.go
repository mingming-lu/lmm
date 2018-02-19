package category

import (
	model "lmm/api/domain/model/category"
	repo "lmm/api/domain/repository/category"
)

func Register(userID, blogID int64, name string) (int64, error) {
	return repo.Add(userID, blogID, name)
}

func Update(userID, blogID int64, name string) error {
	return repo.Update(userID, blogID, name)
}

func FetchByUser(userID int64) ([]model.Category, error) {
	return repo.ByUser(userID)
}

func FetchByBlog(blogID int64) (*model.Category, error) {
	return repo.ByBlog(blogID)
}
