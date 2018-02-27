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

	stmt := d.Must("UPDATE category SET name = ? WHERE user = ? AND category = ?")
	defer stmt.Close()

	res, err := stmt.Exec(name, userID, categoryID)
	rows, err := res.RowsAffected()

	if err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	}

	return nil
}

func ByUser(userID int64) ([]model.Minimal, error) {
	d := db.Default()
	defer d.Close()

	stmt := d.Must("SELECT name FROM category WHERE user = ? GROUP BY name ORDER BY name")
	defer stmt.Close()

	categories := make([]model.Minimal, 0)
	itr, err := stmt.Query(userID)
	if err != nil {
		return categories, err
	}
	for itr.Next() {
		category := model.Minimal{}
		err = itr.Scan(&category.Name)
		if err != nil {
			return categories, nil
		}
		categories = append(categories, category)
	}

	return categories, nil
}
