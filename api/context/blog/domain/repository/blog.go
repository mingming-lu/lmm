package repository

import (
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/blog/domain/model"
	"lmm/api/db"
)

type BlogRepository struct {
	repository.Repository
	Add(blog *model.Blog) error
}

type blogRepo struct {
	repository.Default
}

func NewBlogRepository BlogRepository {
	return new(blogRepo)
}

func (repo *BlogRepository) Add(blog *model.Blog) error {
	db := repo.DB().New()
	defer db.Close()

	stmt := db.MustPrepare(`INSERT INTO blog (id, user, title, text, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(blog.ID, blog.UserID, blog.Title(), blog.Text(), blog.CreatedAt(), blog.UpdatedAt())

	if err != nil {
		return err
	}
	return nil
}
