package blog

import (
	model "lmm/api/domain/model/blog"
	repo "lmm/api/domain/repository/blog"
	"strings"
)

func Post(userID int64, title, text string) (int64, error) {
	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	return repo.Add(userID, title, text)
}

func Fetch(id int64) (*model.Blog, error) {
	return repo.ById(id)
}

func FetchByUser(userID int64) ([]model.Blog, error) {
	return repo.ByUser(userID)
}

func Update(id int64, title, text string) error {
	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	return repo.Update(id, title, text)
}

func Delete(id int64) error {
	return repo.Delete(id)
}
