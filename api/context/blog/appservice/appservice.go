package appservice

import (
	"encoding/json"
	"io"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/service"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
	"lmm/api/utils/strings"
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

type Blog struct {
	ID uint64 `json:"id"`
	BlogContent
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BlogContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BlogListPage struct {
	Blog        []*Blog
	HasNextPage bool
}

func (app *AppService) PostNewBlog(user *account.User, requestBody io.ReadCloser) (uint64, error) {
	content := BlogContent{}
	if err := json.NewDecoder(requestBody).Decode(content); err != nil {
		return 0, err
	}

	blog, err := app.blogService.PostBlog(user.ID(), content.Title, content.Text)
	if err != nil {
		return 0, err
	}

	return blog.ID(), err
}

func (app *AppService) GetBlogListByPage(countStr, pageStr string) (*BlogListPage, error) {
	if countStr == "" {
		countStr = "10"
	}
	count, err := strings.StrToInt(countStr)
	if err != nil {
		return nil, ErrInvalidCount
	}

	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strings.StrToInt(pageStr)
	if err != nil {
		return nil, ErrInvalidPage
	}

	blogList, hasNextPage, err := app.blogService.GetBlogListByPage(count, page)
	blogPage := make([]*Blog, len(blogList))

	for index, blog := range blogList {
		blogPage[index].ID = blog.ID()
		blogPage[index].Title = blog.Title()
		blogPage[index].Text = blog.Text()
		blogPage[index].CreatedAt = blog.CreatedAt().UTC().String()
		blogPage[index].UpdatedAt = blog.UpdatedAt().UTC().String()
	}

	return &BlogListPage{
		Blog:        blogPage,
		HasNextPage: hasNextPage,
	}, nil
}

func (app *AppService) GetBlogByID(blogIDStr string) (*model.Blog, error) {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return nil, service.ErrInvalidBlogID
	}
	return app.blogService.GetBlogByID(blogID)
}

func (app *AppService) GetCategoryOfBlog(blogIDStr string) (*model.Category, error) {
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
	return category, nil
}

func (app *AppService) SetBlogCategory(blogIDStr, categoryName string) error {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return err
	}
	blog, err := app.blogService.GetBlogByID(blogID)
	if err != nil {
		return err
	}
	category, err := app.categoryService.GetCategoryByName(categoryName)
	if err != nil {
		return nil
	}
	return app.blogService.SetBlogCategory(blog, category)
}
