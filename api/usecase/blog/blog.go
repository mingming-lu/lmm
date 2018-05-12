package blog

import (
	model "lmm/api/domain/model/blog"
	repo "lmm/api/domain/repository/blog"
	"strings"

	"github.com/akinaru-lu/errors"
)

func Post(userID uint64, title, text string) (uint64, error) {
	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	return repo.Add(userID, title, text)
}

func FetchByID(id uint64) (*model.Blog, error) {
	return repo.ById(id)
}

func FetchListByUser(userID uint64) ([]model.ListItem, error) {
	return repo.List(userID)
}

func CheckOwnership(userID, blogID uint64) error {
	blog, err := FetchByID(blogID)
	if err != nil {
		return err
	}
	if blog.User != userID {
		return errors.New("User doesn't own the targer blog")
	}
	return nil
}

func Update(id uint64, title, text string) error {
	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	return repo.Update(id, title, text)
}

func Delete(id uint64) error {
	return repo.Delete(id)
}
