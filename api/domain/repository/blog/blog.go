package blog

import (
	"lmm/api/db"
	model "lmm/api/domain/model/blog"
)

func Add(userID uint64, title, text string) (uint64, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO blog (user, title, text) VALUES (?, ?, ?)")
	defer stmt.Close()

	res, err := stmt.Exec(userID, title, text)
	if err != nil {
		return 0, err
	}

	blogID, err := res.LastInsertId()
	return uint64(blogID), err
}

func ById(id uint64) (*model.Blog, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT id, user, title, text, created_at, updated_at FROM blog WHERE id = ?")
	defer stmt.Close()

	blog := model.Blog{}
	err := stmt.QueryRow(id).Scan(
		&blog.ID, &blog.User, &blog.Title, &blog.Text, &blog.CreatedAt, &blog.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func List(userID uint64) ([]model.ListItem, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT id, title, created_at FROM blog WHERE user = ? ORDER BY created_at DESC")
	defer stmt.Close()

	itr, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	blogList := make([]model.ListItem, 0)
	for itr.Next() {
		blog := model.ListItem{}
		err = itr.Scan(&blog.ID, &blog.Title, &blog.CreatedAt)
		if err != nil {
			return blogList, err
		}
		blogList = append(blogList, blog)
	}

	return blogList, nil
}

func Update(id uint64, title, text string) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("UPDATE blog SET title = ?, text = ? WHERE id = ?")
	defer stmt.Close()

	res, err := stmt.Exec(title, text, id)

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoRows
	}

	return nil
}

func Delete(id uint64) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("DELETE FROM blog WHERE id = ?")
	defer stmt.Close()

	res, err := stmt.Exec(id)

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoRows
	}

	return nil
}
