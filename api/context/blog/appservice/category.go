package appservice

import (
	"errors"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/db"
	"lmm/api/utils/strings"
)

var (
	ErrNoSuchCategory = errors.New("no such category")
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

func (app *CategoryApp) UpdateCategory(categoryIDStr, newName string) error {
	categoryID, err := strings.StrToUint64(categoryIDStr)
	if err != nil {
		return ErrNoSuchCategory
	}

	category, err := app.repo.FindByID(categoryID)
	if err != nil {
		return ErrNoSuchBlog
	}

	err = category.UpdateName(newName)
	if err != nil {
		return model.ErrInvalidCategoryName
	}

	err = app.repo.Update(category)
	if err == db.ErrNoChange {
		return ErrNoSuchBlog
	}

	return nil
}
