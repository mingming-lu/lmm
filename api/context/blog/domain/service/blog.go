package service

import (
	"errors"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
	"lmm/api/utils/strings"
)

var (
	ErrBlogTitleDuplicated = errors.New("blog title duplicated")
	ErrEmptyBlogTitle      = errors.New("blog title cannot be empty")
)

type BlogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

func (s *BlogService) PostBlog(userID uint64, title, text string) (*model.Blog, error) {
	if title == "" {
		return nil, ErrEmptyBlogTitle
	}

	blog, err := factory.NewBlog(userID, title, text)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Add(blog); err != nil {
		key, _, ok := storage.CheckErrorDuplicate(err.Error())
		if !ok {
			return nil, err
		}
		if key == "title" {
			return nil, ErrBlogTitleDuplicated
		}
		return nil, err
	}
	return blog, nil
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
	case storage.ErrNoRows:
		return nil, ErrNoSuchBlog
	default:
		return nil, err
	}
}

func (s *BlogService) SetBlogCategory(blog *model.Blog, category *model.Category) error {
	return s.repo.SetBlogCategory(blog, category)
}
