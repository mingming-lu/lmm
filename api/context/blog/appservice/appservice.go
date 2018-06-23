package appservice

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/service"
	"lmm/api/context/blog/repository"
)

type AppService struct {
	blogService     *service.BlogService
	categoryService *service.CategoryService
}

func New(
	blogRepo repository.BlogRepository,
	categoryRepo repository.CategoryRepository) *AppService {
	return &AppService{
		blogService:     service.NewBlogService(blogRepo),
		categoryService: service.NewCategoryService(categoryRepo),
	}
}

func (app *AppService) GetCategoryOfBlog(blogIDStr string) (*model.Category, error) {
	blog, err := app.blogService.GetBlogByID(blogIDStr)
	if err != nil {
		return nil, err
	}
	category, err := app.categoryService.GetCategoryOf(blog)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (app *AppService) SetBlogCategory(blogIDStr, categoryName string) error {
	blog, err := app.blogService.GetBlogByID(blogIDStr)
	if err != nil {
		return err
	}
	category, err := app.categoryService.GetCategoryByName(categoryName)
	if err != nil {
		return nil
	}
	return app.blogService.SetBlogCategory(blog, category)
}
