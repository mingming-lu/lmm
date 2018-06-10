package appservice

import (
	"errors"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/db"
	repoUtil "lmm/api/domain/repository"
	"lmm/api/utils/strings"
)

var (
	ErrBlogTitleDuplicated = errors.New("blog title duplicated")
	ErrNoSuchBlog          = errors.New("no such blog")
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

func (app *AppService) FindBlogByID(idStr string) (*model.Blog, error) {
	id, err := strings.StrToUint64(idStr)
	if err != nil {
		return nil, ErrNoSuchBlog
	}

	blog, err := app.repo.FindByID(id)
	if err != nil {
		if err.Error() == db.ErrNoRows.Error() {
			return nil, ErrNoSuchBlog
		}
		return nil, err
	}
	return blog, nil
}
