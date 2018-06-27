package appservice

import (
	"encoding/json"
	"io"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/domain/service"
	"lmm/api/utils/strings"
)

func (app *AppService) RegisterNewCategory(user *account.User, requestBody io.ReadCloser) (uint64, error) {
	category := Category{}
	if err := json.NewDecoder(requestBody).Decode(&category); err != nil {
		return 0, err
	}

	categoryModel, err := app.categoryService.RegisterCategory(category.Name)
	if err != nil {
		return 0, err
	}

	return categoryModel.ID(), nil
}

func (app *AppService) EditCategory(user *account.User, categoryIDStr string, requestBody io.ReadCloser) error {
	categoryID, err := strings.StrToUint64(categoryIDStr)
	if err != nil {
		return service.ErrInvalidCategoryID
	}

	category := Category{}
	if err := json.NewDecoder(requestBody).Decode(&category); err != nil {
		return err
	}

	cateogryModel, err := app.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}
	if err := cateogryModel.UpdateName(category.Name); err != nil {
		return err
	}

	return app.categoryService.UpdateCategory(cateogryModel)
}

func (app *AppService) GetAllCategories() (*Categories, error) {
	categoryModels, err := app.categoryService.GetAllCategories()
	if err != nil {
		return nil, err
	}
	categories := make([]*Category, len(categoryModels))
	for index, category := range categoryModels {
		categories[index].Name = category.Name()
	}
	return &Categories{
		Categories: categories,
	}, nil
}

func (app *AppService) GetCategoryOfBlog(blogIDStr string) (*Category, error) {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return nil, err
	}
	blog, err := app.blogService.GetBlogByID(blogID)
	if err != nil {
		return nil, err
	}
	category, err := app.categoryService.GetCategoryOf(blog)
	if err != nil {
		return nil, err
	}

	return &Category{
		Name: category.Name(),
	}, nil
}

func (app *AppService) SetBlogCategory(blogIDStr string, requestBody io.ReadCloser) error {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return err
	}

	blogModel, err := app.blogService.GetBlogByID(blogID)
	if err != nil {
		return err
	}

	category := Category{}
	if err := json.NewDecoder(requestBody).Decode(&category); err != nil {
		return err
	}

	categoryModel, err := app.categoryService.GetCategoryByName(category.Name)
	if err != nil {
		return nil
	}
	return app.blogService.SetBlogCategory(blogModel, categoryModel)
}

func (app *AppService) RemoveCategoryByID(idStr string) error {
	id, err := strings.StrToUint64(idStr)
	if err != nil {
		return service.ErrInvalidCategoryID
	}
	return app.categoryService.RemoveCategoryByID(id)
}
