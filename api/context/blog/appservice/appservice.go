package appservice

import (
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
