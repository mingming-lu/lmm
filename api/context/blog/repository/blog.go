package repository

import (
	"lmm/api/context/blog/domain/model"
	sql "lmm/api/db"
	"lmm/api/domain/repository"
	"time"
)

type BlogRepository interface {
	repository.Repository
	Add(blog *model.Blog) error
	Update(blog *model.Blog) error
	FindAll(count, page int) ([]*model.Blog, error)
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
	defer rows.Close()

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

func (repo *blogRepo) FindAllByCategory(category *model.Category, count, page uint) ([]*model.Blog, error) {
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`
		SELECT b.id, b.user, b.title, b.text, b.created_at, b.updated_at
		FROM blog_category AS bc
		INNER JOIN blog AS b ON b.id = bc.blog
		WHERE bc.category = ?
		ORDER BY b.created_at
		LIMIT ?
		OFFSET ?
	`)
	defer stmt.Close()

	rows, err := stmt.Query(category.ID(), count+1, (page-1)*count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	blogList := make([]*model.Blog, 0)
	var (
		blogID        uint64
		blogWriter    uint64
		blogTitle     string
		blogText      string
		blogCreatedAt time.Time
		blogUpdatedAt time.Time
	)

	for rows.Next() {
		if err := rows.Scan(
			&blogID, &blogWriter, &blogTitle, &blogText, &blogCreatedAt, &blogUpdatedAt,
		); err != nil {
			return nil, err
		}
		blog := model.NewBlog(blogID, blogWriter, blogTitle, blogText, blogCreatedAt, blogUpdatedAt)
		blogList = append(blogList, blog)
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

	res, err := stmt.Exec(blog.Title(), blog.Text(), blog.ID(), blog.UserID())
	if err != nil {
		return err
	}
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
