package repository

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/storage"
)

type CategoryRepository interface {
	Add(category *model.Category) error
	Update(categoryRepo *model.Category) error
	FindAll() ([]*model.Category, error)
	FindByID(id uint64) (*model.Category, error)
	FindByName(name string) (*model.Category, error)
	FindByBlog(blog *model.Blog) (*model.Category, error)
	Remove(category *model.Category) error
}

type categoryRepo struct {
	db *storage.DB
}

func NewCategoryRepository(db *storage.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (repo *categoryRepo) Add(category *model.Category) error {
	stmt := repo.db.MustPrepare(`INSERT INTO category (id, name) VALUES (?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(category.ID(), category.Name())
	return err
}

func (repo *categoryRepo) Update(category *model.Category) error {
	stmt := repo.db.MustPrepare(`UPDATE category SET name = ? WHERE id = ?`)
	defer stmt.Close()

	res, err := stmt.Exec(category.Name(), category.ID())
	rowsAffeted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffeted == 0 {
		// wether or not the blog with id exists, return no change
		return storage.ErrNoChange
	}
	return nil
}

func (repo *categoryRepo) FindAll() ([]*model.Category, error) {
	stmt := repo.db.MustPrepare(`SELECT id, name FROM category ORDER BY name`)
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		categoryID   uint64
		categoryName string
	)

	categories := make([]*model.Category, 0)
	for rows.Next() {
		err = rows.Scan(&categoryID, &categoryName)
		if err != nil {
			return nil, err
		}
		category, err := model.NewCategory(categoryID, categoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (repo *categoryRepo) FindByID(id uint64) (*model.Category, error) {
	stmt := repo.db.MustPrepare(`SELECT id, name FROM category WHERE id = ?`)
	defer stmt.Close()

	var (
		categoryID   uint64
		categoryName string
	)

	err := stmt.QueryRow(id).Scan(&categoryID, &categoryName)
	if err != nil {
		return nil, err
	}

	return model.NewCategory(categoryID, categoryName)
}

func (repo *categoryRepo) FindByName(name string) (*model.Category, error) {
	stmt := repo.db.MustPrepare(`SELECT id, name FROM category WHERE name = ?`)
	defer stmt.Close()

	var (
		categoryID   uint64
		categoryName string
	)

	err := stmt.QueryRow(name).Scan(&categoryID, &categoryName)
	if err != nil {
		return nil, err
	}

	return model.NewCategory(categoryID, categoryName)
}

func (repo *categoryRepo) FindByBlog(blog *model.Blog) (*model.Category, error) {
	stmt := repo.db.MustPrepare(`
		SELECT c.id, c.name
		FROM blog_category AS bc
		INNER JOIN category AS c ON c.id = bc.category
		WHERE bc.blog = ?
	`)
	defer stmt.Close()

	var (
		categoryID   uint64
		categoryName string
	)

	err := stmt.QueryRow(blog.ID()).Scan(&categoryID, &categoryName)
	if err != nil {
		return nil, err
	}

	return model.NewCategory(categoryID, categoryName)
}

func (repo *categoryRepo) Remove(category *model.Category) error {
	stmt := repo.db.MustPrepare(`DELETE FROM category WHERE id = ? AND name = ?`)
	defer stmt.Close()

	res, err := stmt.Exec(category.ID(), category.Name())
	if err != nil {
		return err
	}

	rowsAffeted, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffeted == 0 {
		return storage.ErrNoRows
	}

	return nil
}
