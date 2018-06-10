package repository

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/domain/repository"
	"time"
)

type BlogRepository interface {
	repository.Repository
	Add(blog *model.Blog) error
	Update(blog *model.Blog) error
	FindAll(int, int) ([]*model.Blog, error)
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

func (repo *blogRepo) FindAll(count, page int) ([]*model.Blog, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`
		SELECT id, user, title, text, created_at, updated_at
		FROM blog
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
	`)
	defer stmt.Close()

	rows, err := stmt.Query(count+1, (page-1)*count)
	if err != nil {
		return nil, nil
	}

	var (
		blogID        uint64
		blogWriter    uint64
		blogTitle     string
		blogText      string
		blogCreatedAt time.Time
		blogUpdated   time.Time
	)

	blogList := make([]*model.Blog, 0)

	for rows.Next() {
		err = rows.Scan(&blogID, &blogWriter, &blogTitle, &blogText, &blogCreatedAt, &blogUpdated)
		if err != nil {
			return nil, err
		}
		blogList = append(blogList, model.NewBlog(
			blogID, blogWriter, blogTitle, blogText, blogCreatedAt, blogUpdated,
		))
	}

	return blogList, nil
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

func (repo *blogRepo) Update(blog *model.Blog) error {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`UPDATE blog SET title = ?, text = ? WHERE id = ? and user = ?`)
	defer stmt.Close()

	_, err := stmt.Exec(blog.Title(), blog.Text(), blog.ID(), blog.UserID())
	if err != nil {
		return err
	}
	return nil
}
