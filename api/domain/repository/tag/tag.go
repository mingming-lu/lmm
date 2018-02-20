package tag

import (
	"lmm/api/db"
	model "lmm/api/domain/model/tag"

	"github.com/akinaru-lu/errors"
)

func Add(userID, blogID int64, tags []model.Minimal) error {
	d := db.UseDefault()
	defer d.Close()

	stmt := d.Must("INSERT INTO tag (user, blog, name) VALUES (?, ?, ?)")

	tx, err := d.Begin()
	stmt = tx.Stmt(stmt)
	defer stmt.Close()

	rowsAffected := int64(0)
	for _, tag := range tags {
		res, err := stmt.Exec(userID, blogID, tag.Name)
		if err != nil {
			break
		}
		rows, err := res.RowsAffected()
		if err != nil {
			break
		}
		rowsAffected += rows
	}

	if err != nil {
		return errors.Wrap(err, errors.Wrap(tx.Rollback(), "").Error())
	}
	if rowsAffected == 0 {
		return db.ErrNoChange
	}
	if rowsAffected != int64(len(tags)) {
		return errors.Wrap(err, errors.Wrap(tx.Rollback(), "Rows inserted not equal to length of input array").Error())
	}
	return tx.Commit()
}

func Update(userID, blogID, tagID int64, name string) error {
	d := db.UseDefault()
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

func Delete(userID, blogID, tagID int64) error {
	d := db.UseDefault()
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
