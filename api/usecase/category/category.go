package category

import (
	"lmm/api/domain/model/blog"
	model "lmm/api/domain/model/category"
	repo "lmm/api/domain/repository/category"
	"strings"

	"github.com/akinaru-lu/errors"
)

func Register(userID int64, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("Empty name")
	}
	return repo.Add(userID, name)
}

func Update(userID, categoryID int64, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("Empty name")
	}
	return repo.Update(userID, categoryID, name)
}

func FetchByUser(userID int64) ([]model.Category, error) {
	return repo.ByUser(userID)
}

func FetchByBlog(blogID int64) (*model.Category, error) {
	return repo.ByBlog(blogID)
}

func FetchAllBlog(categoryID int64) ([]blog.ListItem, error) {
	return repo.AllBlogByID(categoryID)
}

func Delete(userID, categoryID int64) error {
	blogList, err := FetchAllBlog(categoryID)
	if err != nil {
		return err
	}
	if blogList != nil && len(blogList) != 0 {
		return errors.New("There are still blog in this category")
	}
	return repo.Delete(userID, categoryID)
}
