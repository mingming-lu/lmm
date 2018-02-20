package tag

import (
	"lmm/api/db"
	model "lmm/api/domain/model/tag"

	"github.com/akinaru-lu/errors"
)

func Add(userID, blogID int64, tags []model.Minimal) error {
	return nil
}

func Update(userID, blogID, tagID int64, name string) error {
	return nil
}

func ByID(id int64) (*model.Tag, error) {
	tags, err := bySingle("id", id)
	if err != nil {
		return nil, err
	}
	if len(tags) > 1 {
		return nil, errors.Newf("Unexpected length of tags: %d", len(tags))
	}
	return &tags[0], nil
}

func ByUser(userID int64) ([]model.Tag, error) {
	return bySingle("user", userID)
}

func ByBlog(blogID int64) ([]model.Tag, error) {
	return bySingle("blog", blogID)
}

func bySingle(field string, value interface{}) ([]model.Tag, error) {
	d := db.UseDefault()
	defer d.Close()

	stmt := d.Mustf("SELECT id, user, blog, name FROM tag WHERE %s = ?", field)

	tags := make([]model.Tag, 0)

	itr, err := stmt.Query(value)
	if err != nil {
		return tags, err
	}

	for itr.Next() {
		tag := model.Tag{}
		err = itr.Scan(&tag.ID, &tag.User, &tag.Blog, &tag.Name)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func Delete(userID, blogID, tagID int64) error {
	return nil
}
