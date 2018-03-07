package blog

import (
	model "lmm/api/domain/model/blog"
	repo "lmm/api/domain/repository/blog"
	"strings"

	"github.com/akinaru-lu/errors"
)

func Post(userID int64, title, text string) (int64, error) {
	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	return repo.Add(userID, title, text)
}

func FetchByID(id int64) (*model.Blog, error) {
	return repo.ById(id)
}

func FetchByUser(userID int64) ([]model.Blog, error) {
	return repo.ByUser(userID)
}

func FetchListByUser(userID int64) ([]model.ListItem, error) {
	return repo.List(userID)
}

func CheckOwnership(userID, blogID int64) error {
	blog, err := FetchByID(blogID)
	if err != nil {
		return err
	}
	if blog.User != userID {
		return errors.New("User doesn't own the targer blog")
	}
	return nil
}

func Update(id int64, title, text string) error {
	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	return repo.Update(id, title, text)
}

func Delete(id int64) error {
	return repo.Delete(id)
}
