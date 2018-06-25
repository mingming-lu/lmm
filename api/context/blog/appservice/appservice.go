package appservice

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/service"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
)

type AppService struct {
	blogService     *service.BlogService
	categoryService *service.CategoryService
}

func New(db *storage.DB) *AppService {
	return &AppService{
		blogService:     service.NewBlogService(repository.NewBlogRepository(db)),
		categoryService: service.NewCategoryService(repository.NewCategoryRepository(db)),
	}
}

func (app *AppService) PostNewBlog(user *account.User, title, text string) (uint64, error) {
	blog, err := app.blogService.PostBlog(user.ID(), title, text)
	if err != nil {
		return 0, err
	}

	return blog.ID(), err
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
