package appservice

import (
	"errors"
	"lmm/api/context/blog/repository"
	"lmm/api/storage"
	"lmm/api/utils/strings"
)

var (
	ErrBlogTitleDuplicated = errors.New("blog title duplicated")
	ErrNoSuchBlog          = errors.New("no such blog")
	ErrEmptyBlogTitle      = errors.New("blog title cannot be empty")
	ErrBlogNoChange        = errors.New("blog no change")
	ErrInvalidCount        = errors.New("invalid count")
	ErrInvalidPage         = errors.New("invalid page")
	ErrNoPermission        = errors.New("no permission")
)

type BlogApp struct {
	repo repository.BlogRepository
}

func NewBlogApp(repo repository.BlogRepository) *BlogApp {
	return &BlogApp{repo: repo}
}

func (app *BlogApp) EditBlog(userID uint64, blogIDStr, title, text string) error {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return ErrNoSuchBlog
	}

	blog, err := app.repo.FindByID(blogID)
	if err != nil {
		return ErrNoSuchBlog
	}

	if blog.UserID() != userID {
		return ErrNoPermission
	}

	lastUpdated := blog.UpdatedAt()

	// TODO move validation to model
	if title == "" {
		return ErrEmptyBlogTitle
	}

	blog.UpdateTitle(title)
	blog.UpdateText(text)

	if blog.UpdatedAt().Equal(lastUpdated) {
		return ErrBlogNoChange
	}

	err = app.repo.Update(blog)
	if err == storage.ErrNoChange {
		return ErrNoSuchBlog
	}

	return err // unknown error or nil
}
