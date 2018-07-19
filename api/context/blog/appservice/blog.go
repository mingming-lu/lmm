package appservice

import (
	"encoding/json"
	"io"
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/domain"
	"lmm/api/utils/strings"
)

func (app *AppService) PostNewBlog(user *account.User, requestBody io.ReadCloser) (uint64, error) {
	content := BlogContent{}
	if err := json.NewDecoder(requestBody).Decode(&content); err != nil {
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
		return nil, domain.ErrInvalidCount
	}

	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strings.StrToInt(pageStr)
	if err != nil {
		return nil, domain.ErrInvalidPage
	}

	blogList, nextPage, err := app.blogService.GetBlogListByPage(count, page)
	if err != nil {
		return nil, err
	}

	blogPage := make([]*Blog, len(blogList))

	for index, blog := range blogList {
		data := &Blog{
			ID: blog.ID(),
			BlogContent: BlogContent{
				Title: blog.Title(),
				Text:  blog.Text(),
			},
			CreatedAt: blog.CreatedAt().UTC().String(),
			UpdatedAt: blog.UpdatedAt().UTC().String(),
		}
		blogPage[index] = data
	}

	return &BlogListPage{
		Blog:     blogPage,
		NextPage: nextPage,
	}, nil
}

func (app *AppService) GetBlogByID(blogIDStr string) (*Blog, error) {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return nil, domain.ErrInvalidBlogID
	}
	blog, err := app.blogService.GetBlogByID(blogID)
	if err != nil {
		return nil, err
	}
	return &Blog{
		ID: blog.ID(),
		BlogContent: BlogContent{
			Title: blog.Title(),
			Text:  blog.Text(),
		},
		CreatedAt: blog.CreatedAt().UTC().String(),
		UpdatedAt: blog.UpdatedAt().UTC().String(),
	}, nil
}

func (app *AppService) EditBlog(user *account.User, blogIDStr string, requestBody io.ReadCloser) error {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return domain.ErrNoSuchBlog
	}

	content := BlogContent{}
	if err := json.NewDecoder(requestBody).Decode(&content); err != nil {
		return err
	}

	return app.blogService.EditBlog(user.ID(), blogID, content.Title, content.Text)
}
