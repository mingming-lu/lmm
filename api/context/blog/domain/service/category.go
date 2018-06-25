package service

import (
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) RegisterCategory(name string) (*model.Category, error) {
	category, err := factory.NewCategory(name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Add(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetAllCategories() ([]*model.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetCategoryByID(id uint64) (*model.Category, error) {
	category, err := s.repo.FindByID(id)
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

func (s *CategoryService) UpdateCategory(category *model.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) RemoveCategoryByID(id uint64) error {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return ErrNoSuchCategory
	}

	return s.repo.Remove(category)
}
