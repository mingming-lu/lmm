package tag

import (
	model "lmm/api/domain/model/tag"
	repo "lmm/api/domain/repository/tag"
	"lmm/api/usecase/blog"

	"github.com/akinaru-lu/errors"
)

func Add(userID, blogID uint64, tagName string) (uint64, error) {
	blog, err := blog.FetchByID(blogID)
	if err != nil {
		return 0, errors.Wrapf(err, "No such blog: %d", blogID)
	}

	return repo.Add(userID, blog.ID, tagName)
}

func Update(userID, blogID, tagID uint64, name string) error {
	tag, err := repo.ByID(tagID)
	if err != nil {
		return errors.Wrapf(err, "No such tag: %d", tagID)
	}

	blog, err := blog.FetchByID(blogID)
	if err != nil {
		return errors.Wrapf(err, "No such blog: %d", blogID)
	}

	return repo.Update(userID, blog.ID, tag.ID, name)
}

func FetchByUser(userID uint64) ([]model.Tag, error) {
	return repo.ByUser(userID)
}

func FetchByBlog(blogID uint64) ([]model.Tag, error) {
	return repo.ByBlog(blogID)
}

func Delete(userID, blogID, tagID uint64) error {
	return repo.Delete(userID, blogID, tagID)
}
