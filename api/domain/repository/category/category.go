package category

import (
	"lmm/api/db"
	model "lmm/api/domain/model/category"
)

func Add(userID int64, name string) (int64, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("INSERT INTO category (user, name) VALUES (?, ?)")
	defer stmt.Close()

	res, err := stmt.Exec(userID, name)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Update(userID, categoryID int64, name string) error {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("UPDATE category SET name = ? WHERE user = ? AND id = ?")
	defer stmt.Close()

	res, err := stmt.Exec(name, userID, categoryID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	}

	return nil
}

func ByUser(userID int64) ([]model.Category, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT id, user, name FROM category WHERE user = ? ORDER BY name")
	defer stmt.Close()

	categories := make([]model.Category, 0)
	itr, err := stmt.Query(userID)
	if err != nil {
		return categories, err
	}
	for itr.Next() {
		category := model.Category{}
		err = itr.Scan(&category.ID, &category.User, &category.Name)
		if err != nil {
			return categories, nil
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func ByBlog(blogID int64) (*model.Category, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT c.id, c.user, c.name FROM blog_category AS bc INNER JOIN category AS c ON c.id = bc.category WHERE bc.blog = ?")
	defer stmt.Close()

	category := model.Category{}
	err := stmt.QueryRow(blogID).Scan(&category.ID, &category.User, &category.Name)

	return &category, err
}
