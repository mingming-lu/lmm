package category

import (
	"lmm/api/db"
	model "lmm/api/domain/model/category"
)

func Add(userID, blogID int64, name string) (int64, error) {
	d := db.UseDefault()
	defer d.Close()

	stmt := d.Must("INSERT INTO category (user, blog, name) VALUES (?, ?, ?)")
	defer stmt.Close()

	res, err := stmt.Exec(userID, blogID, name)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Update(userID, blogID int64, name string) error {
	d := db.UseDefault()
	defer d.Close()

	stmt := d.Must("UPDATE category SET name = ? WHERE user = ? AND blog = ?")
	defer stmt.Close()

	res, err := stmt.Exec(name, userID, blogID)
	rows, err := res.RowsAffected()

	if err != nil {
		return err
	} else if rows == 0 {
		return db.ErrNoChange
	}

	return nil
}

func ByUser(userID int64) ([]model.Category, error) {
	d := db.UseDefault()
	defer d.Close()

	stmt := d.Must("SELECT MIN(id), MIN(user), MIN(blog), name FROM category WHERE user = ? GROUP BY name ORDER BY name")
	defer stmt.Close()

	categories := make([]model.Category, 0)
	itr, err := stmt.Query(userID)
	if err != nil {
		return categories, err
	}
	for itr.Next() {
		category := model.Category{}
		err = itr.Scan(&category.ID, &category.User, &category.Blog, &category.Name)
		if err != nil {
			return categories, nil
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func ByBlog(blogID int64) (*model.Category, error) {
	d := db.UseDefault()
	defer d.Close()

	stmt := d.Must("SELECT id, user, blog, name FROM category WHERE blog = ?")
	defer stmt.Close()

	category := model.Category{}
	err := stmt.QueryRow(blogID).Scan(&category.ID, &category.User, &category.Blog, &category.Name)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
