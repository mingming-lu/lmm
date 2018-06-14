package appservice

import (
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/repository"
)

type CategoryApp struct {
	repo repository.CategoryRepository
}

func NewCategoryApp(repo repository.CategoryRepository) *CategoryApp {
	return &CategoryApp{repo: repo}
}

func (app *CategoryApp) AddNewCategory(name string) (uint64, error) {
	category, err := factory.NewCategory(name)
	if err != nil {
		return 0, err
	}

	err = app.repo.Add(category)
	if err != nil {
		return 0, err
	}
	return category.ID(), nil
}
