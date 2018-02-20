package tag

import (
	model "lmm/api/domain/model/tag"
	repo "lmm/api/domain/repository/tag"
	"lmm/api/usecase/blog"

	"github.com/akinaru-lu/errors"
)

func Add(userID, blogID int64, tags []model.Minimal) error {
	blog, err := blog.Fetch(blogID)
	if err != nil {
		return errors.Wrapf(err, "No such blog: %d", blogID)
	}

	return repo.Add(userID, blog.ID, tags)
}

func Update(userID, blogID, tagID int64, name string) error {
	tag, err := repo.ByID(tagID)
	if err != nil {
		return errors.Wrapf(err, "No such tag: %d", tagID)
	}

	blog, err := blog.Fetch(blogID)
	if err != nil {
		return errors.Wrapf(err, "No such blog: %d", blogID)
	}

	return repo.Update(userID, blog.ID, tag.ID, name)
}

func FetchByUser(userID int64) ([]model.Tag, error) {
	return repo.ByUser(userID)
}

func FetchByBlog(blogID int64) ([]model.Tag, error) {
	return repo.ByBlog(blogID)
}

func Delete(userID, blogID, tagID int64) error {
	return repo.Delete(userID, blogID, tagID)
}
