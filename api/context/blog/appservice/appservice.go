package appservice

import (
	"errors"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/repository"
	repoUtil "lmm/api/domain/repository"
)

var (
	ErrBlogTitleDuplicated = errors.New("blog title duplicated")
)

type AppService struct {
	repo repository.BlogRepository
}

func New(repo repository.BlogRepository) *AppService {
	return &AppService{repo: repo}
}

func (app *AppService) PostNewBlog(userID uint64, title, text string) (uint64, error) {
	blog, err := factory.NewBlog(userID, title, text)
	if err != nil {
		return uint64(0), err
	}
	err = app.repo.Add(blog)
	if err != nil {
		key, _, ok := repoUtil.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "title" {
			return 0, ErrBlogTitleDuplicated
		}
		return 0, err
	}
	return blog.ID(), err
}
