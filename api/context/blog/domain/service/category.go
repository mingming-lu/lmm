package service

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
	"lmm/api/utils/strings"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategoryByID(idStr string) (*model.Category, error) {
	categoryID, err := strings.StrToUint64(idStr)
	if err != nil {
		return nil, ErrInvalidCategoryID
	}
	category, err := s.repo.FindByID(categoryID)
	switch err {
	case nil:
		return category, nil
	case storage.ErrNoRows:
		return nil, ErrNoSuchCategory
	default:
		return nil, err
	}
}

func (s *CategoryService) GetCategoryByName(name string) (*model.Category, error) {
	category, err := s.repo.FindByName(name)
	switch err {
	case nil:
		return category, nil
	case storage.ErrNoRows:
		return nil, ErrNoSuchCategory
	default:
		return nil, err
	}
}

func (s *CategoryService) GetCategoryOf(blog *model.Blog) (*model.Category, error) {
	category, err := s.repo.FindByBlog(blog)
	switch err {
	case nil:
		return category, nil
	case storage.ErrNoRows:
		return nil, ErrCategoryNotSet
	default:
		return nil, err
	}
}
