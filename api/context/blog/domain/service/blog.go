package service

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/repository"
	"lmm/api/db"
	"lmm/api/utils/strings"
)

type BlogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

func (s *BlogService) GetBlogByID(blogIDStr string) (*model.Blog, error) {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return nil, ErrInvalidBlogID
	}
	blog, err := s.repo.FindByID(blogID)
	switch err {
	case nil:
		return blog, nil
	case db.ErrNoRows:
		return nil, ErrNoSuchBlog
	default:
		return nil, err
	}
}

func (s *BlogService) SetBlogCategory(blog *model.Blog, category *model.Category) error {
	return s.repo.SetBlogCategory(blog, category)
}
