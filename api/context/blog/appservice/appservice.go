package appservice

import (
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/repository"
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
	return blog.ID(), app.repo.Add(blog)
}
