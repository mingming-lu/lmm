package service

import (
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/storage"
)

type BlogService struct {
	repo repository.BlogRepository
}

func NewBlogService(repo repository.BlogRepository) *BlogService {
	return &BlogService{repo: repo}
}

func (s *BlogService) PostBlog(userID uint64, title, text string) (*model.Blog, error) {
	if title == "" {
		return nil, ErrEmptyBlogTitle
	}

	blog, err := factory.NewBlog(userID, title, text)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Add(blog); err != nil {
		key, _, ok := storage.CheckErrorDuplicate(err)
		if !ok {
			return nil, err
		}
		if key == "title" {
			return nil, ErrBlogTitleDuplicated
		}
		return nil, err
	}
	return blog, nil
}

func (s *BlogService) GetBlogListByPage(count, page int) ([]*model.Blog, int, error) {
	return s.repo.FindAll(count, page)
}

func (s *BlogService) GetBlogByID(id uint64) (*model.Blog, error) {
	blog, err := s.repo.FindByID(id)
	switch err {
	case nil:
		return blog, nil
	case storage.ErrNoRows:
		return nil, ErrNoSuchBlog
	default:
		return nil, err
	}
}

func (s *BlogService) EditBlog(userID, blogID uint64, title, text string) error {
	blog, err := s.repo.FindByID(blogID)
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

	err = s.repo.Update(blog)
	if err == storage.ErrNoChange {
		return ErrNoSuchBlog
	}

	return err
}

func (s *BlogService) SetBlogCategory(blog *model.Blog, category *model.Category) error {
	return s.repo.SetBlogCategory(blog, category)
}
