package appservice

import (
	"lmm/api/context/blog/domain/repository"
	"lmm/api/context/blog/domain/service"
)

type AppService struct {
	blogService     *service.BlogService
	categoryService *service.CategoryService
}

func New(
	blogRepo repository.BlogRepository,
	categoryRepo repository.CategoryRepository,
) *AppService {
	return &AppService{
		blogService:     service.NewBlogService(blogRepo),
		categoryService: service.NewCategoryService(categoryRepo),
	}
}
