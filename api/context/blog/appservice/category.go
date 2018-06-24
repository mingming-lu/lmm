package appservice

import (
	"errors"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
	"lmm/api/utils/strings"
)

var (
	ErrCategoryNoChanged     = errors.New("category no changed")
	ErrDuplicateCategoryName = errors.New("duplicate category name")
	ErrInvalidCategoryName   = errors.New("invalid category name")
	ErrNoSuchCategory        = errors.New("no such category")
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
		return 0, ErrInvalidCategoryName
	}

	err = app.repo.Add(category)
	if err != nil {
		key, _, ok := storage.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "name" {
			return 0, ErrDuplicateCategoryName
		}
		return 0, err
	}
	return category.ID(), nil
}

func (app *CategoryApp) UpdateCategoryName(categoryIDStr, newName string) error {
	categoryID, err := strings.StrToUint64(categoryIDStr)
	if err != nil {
		return ErrNoSuchCategory
	}

	category, err := app.repo.FindByID(categoryID)
	if err != nil {
		return ErrNoSuchCategory
	}

	err = category.UpdateName(newName)
	if err != nil {
		return ErrInvalidCategoryName
	}

	err = app.repo.Update(category)
	if err == storage.ErrNoChange {
		return ErrCategoryNoChanged
	}

	return nil
}

func (app *CategoryApp) FindAllCategories() ([]*model.Category, error) {
	categories, err := app.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (app *CategoryApp) Remove(idStr string) error {
	id, err := strings.StrToUint64(idStr)
	if err != nil {
		return ErrNoSuchCategory
	}

	category, err := app.repo.FindByID(id)
	if err != nil {
		return ErrNoSuchCategory
	}

	switch app.repo.Remove(category) {
	case nil:
		return nil
	case storage.ErrNoRows:
		return ErrNoSuchCategory
	default:
		return err
	}
}
