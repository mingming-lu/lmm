package tag

import (
	"lmm/api/db"
	model "lmm/api/domain/model/tag"

	"github.com/akinaru-lu/errors"
)

func Add(userID, blogID uint64, tagName string) (uint64, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO tag (user, blog, name) VALUES (?, ?, ?)")
	defer stmt.Close()

	res, err := stmt.Exec(userID, blogID, tagName)
	if err != nil {
		return 0, err
	}
	tagID, err := res.LastInsertId()
	return uint64(tagID), err
}

func Update(userID, blogID, tagID uint64, name string) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("UPDATE tag SET name = ? WHERE id = ? AND user = ? AND blog = ?")
	defer stmt.Close()

	res, err := stmt.Exec(name, tagID, userID, blogID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return db.ErrNoChange
	}

	return nil
}

func ByID(id uint64) (*model.Tag, error) {
	tags, err := bySingle("id", id)
	if err != nil {
		return nil, err
	}
	if len(tags) > 1 {
		return nil, errors.Newf("Unexpected length of tags: %d", len(tags))
	}
	return &tags[0], nil
}

func ByUser(userID uint64) ([]model.Tag, error) {
	return bySingle("user", userID)
}

func ByBlog(blogID uint64) ([]model.Tag, error) {
	return bySingle("blog", blogID)
}

func bySingle(field string, value interface{}) ([]model.Tag, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Mustf("SELECT MIN(id), MIN(user), MIN(blog), name FROM tag WHERE %s = ? GROUP BY name ORDER BY name", field)
	defer stmt.Close()

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

func Delete(userID, blogID, tagID uint64) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("DELETE FROM tag WHERE id = ? AND user = ? AND blog = ?")
	defer stmt.Close()

	res, err := stmt.Exec(tagID, userID, blogID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return db.ErrNoChange
	}

	return nil
}
