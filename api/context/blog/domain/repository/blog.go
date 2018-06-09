package repository

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/domain/repository"
	"time"
)

type BlogRepository interface {
	repository.Repository
	Add(blog *model.Blog) error
	FindByID(id uint64) (*model.Blog, error)
}

type blogRepo struct {
	repository.Default
}

func NewBlogRepository() BlogRepository {
	return new(blogRepo)
}

func (repo *blogRepo) Add(blog *model.Blog) error {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`INSERT INTO blog (id, user, title, text, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(blog.ID(), blog.UserID(), blog.Title(), blog.Text(), blog.CreatedAt().UTC(), blog.UpdatedAt().UTC())

	if err != nil {
		return err
	}
	return nil
}

func (repo *blogRepo) FindByID(id uint64) (*model.Blog, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT id, user, title, text, created_at, updated_at FROM blog WHERE id = ?`)
	defer stmt.Close()

	var (
		blogID       uint64
		blogWriter   uint64
		blogTitle    string
		blogText     string
		blogCreateAt time.Time
		blogUpdaedAt time.Time
	)

	err := stmt.QueryRow(id).Scan(&blogID, &blogWriter, &blogTitle, &blogText, &blogCreateAt, &blogUpdaedAt)

	if err != nil {
		return nil, err
	}
	return model.NewBlog(blogID, blogWriter, blogTitle, blogText, blogCreateAt, blogUpdaedAt), nil
}
