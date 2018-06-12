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
	ErrBlogNoChange        = errors.New("blog no change")
	ErrInvalidCount        = errors.New("invalid count")
	ErrInvalidPage         = errors.New("invalid page")
	ErrNoPermission        = errors.New("no permission")
)

type AppService struct {
	repo repository.BlogRepository
}

func New(repo repository.BlogRepository) *AppService {
	return &AppService{repo: repo}
}

func (app *AppService) PostNewBlog(userID uint64, title, text string) (uint64, error) {
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

func (app *AppService) FindAllBlog(countStr, pageStr string) ([]*model.Blog, int, bool, error) {
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

func (app *AppService) FindBlogByID(idStr string) (*model.Blog, error) {
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

func (app *AppService) EditBlog(userID uint64, blogIDStr, title, text string) error {
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
