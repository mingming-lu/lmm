package repository

import (
	"lmm/api/context/blog/domain/model"
	sql "lmm/api/db"
	"lmm/api/domain/repository"
)

type CategoryRepository interface {
	repository.Repository
	Add(category *model.Category) error
	Update(categoryRepo *model.Category) error
	FindAll() ([]*model.Category, error)
	FindByID(id uint64) (*model.Category, error)
}

type categoryRepo struct {
	repository.Default
}

func NewCategoryRepository() CategoryRepository {
	return new(categoryRepo)
}

func (repo *categoryRepo) Add(category *model.Category) error {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`INSERT INTO category (id, name) VALUES (?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(category.ID(), category.Name())
	return err
}

func (repo *categoryRepo) Update(category *model.Category) error {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`UPDATE category SET name = ? WHERE id = ?`)
	defer stmt.Close()

	res, err := stmt.Exec(category.Name(), category.ID())
	rowsAffeted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffeted == 0 {
		// wether or not the blog with id exists, return no change
		return sql.ErrNoChange
	}
	return nil
}

func (repo *categoryRepo) FindAll() ([]*model.Category, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name FROM category ORDER BY name`)
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

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
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, name FROM category WHERE id = ?`)
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
