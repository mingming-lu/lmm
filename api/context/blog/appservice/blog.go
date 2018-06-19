package appservice

import (
	"errors"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/db"
	repoUtil "lmm/api/domain/repository"
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

func (app *BlogApp) PostNewBlog(userID uint64, title, text string) (uint64, error) {
	if title == "" {
		return 0, ErrEmptyBlogTitle
	}
	blog, err := factory.NewBlog(userID, title, text)
	if err != nil {
		return uint64(0), err
	}
	err = app.repo.Add(blog)
	if err != nil {
		key, _, ok := repoUtil.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "title" {
			return 0, ErrBlogTitleDuplicated
		}
		return 0, err
	}
	return blog.ID(), err
}

func (app *BlogApp) FindAllBlog(countStr, pageStr string) ([]*model.Blog, int, bool, error) {
	var (
		count int
		page  int
		err   error
	)
	if countStr == "" {
		count = 10
	} else {
		count, err = strings.StrToInt(countStr)
		if err != nil {
			return nil, 0, false, ErrInvalidCount
		}
	}
	if pageStr == "" {
		page = 1
	} else {
		page, err = strings.StrToInt(pageStr)
		if err != nil {
			return nil, 0, false, ErrInvalidPage
		}
	}

	blogList, err := app.repo.FindAll(count, page)
	if err != nil {
		return nil, 0, false, err
	}

	if len(blogList) <= count {
		return blogList, page, false, nil
	}
	return blogList[:count], page, true, nil
}

func (app *BlogApp) FindBlogByID(idStr string) (*model.Blog, error) {
	id, err := strings.StrToUint64(idStr)
	if err != nil {
		return nil, ErrNoSuchBlog
	}

	blog, err := app.repo.FindByID(id)
	if err != nil {
		if err.Error() == db.ErrNoRows.Error() {
			return nil, ErrNoSuchBlog
		}
		return nil, err
	}
	return blog, nil
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
	if err == db.ErrNoChange {
		return ErrNoSuchBlog
	}

	return err // unknown error or nil
}
