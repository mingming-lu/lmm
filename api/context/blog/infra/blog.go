package infra

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/storage"
	"time"
)

type BlogStorage struct {
	db *storage.DB
}

func NewBlogStorage(db *storage.DB) *BlogStorage {
	return &BlogStorage{db: db}
}

func (s *BlogStorage) Add(blog *model.Blog) error {
	stmt := s.db.MustPrepare(`INSERT INTO blog (id, user, title, text, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(blog.ID(), blog.UserID(), blog.Title(), blog.Text(), blog.CreatedAt().UTC(), blog.UpdatedAt().UTC())

	if err != nil {
		return err
	}
	return nil
}

func (s *BlogStorage) FindAll(count, page int) ([]*model.Blog, int, error) {
	stmt := s.db.MustPrepare(`
		SELECT id, user, title, text, created_at, updated_at
		FROM blog
		ORDER BY created_at DESC
		LIMIT ?
		OFFSET ?
	`)
	defer stmt.Close()

	rows, err := stmt.Query(count+1, (page-1)*count)
	if err != nil {
		return nil, -1, err
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
			return nil, -1, err
		}
		blogList = append(blogList, model.NewBlog(
			blogID, blogWriter, blogTitle, blogText, blogCreatedAt, blogUpdated,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, -1, err
	}

	nextPage := -1
	if len(blogList) > count {
		blogList = blogList[:count]
		nextPage = page + 1
	}

	return blogList, nextPage, nil
}

func (s *BlogStorage) FindAllByCategory(category *model.Category, count, page int) ([]*model.Blog, int, error) {
	stmt := s.db.MustPrepare(`
		SELECT b.id, b.user, b.title, b.text, b.created_at, b.updated_at
		FROM blog_category AS bc
		INNER JOIN blog AS b ON b.id = bc.blog
		WHERE bc.category = ?
		ORDER BY bc.blog DESC
		LIMIT ?
		OFFSET ?
	`)
	defer stmt.Close()

	rows, err := stmt.Query(category.ID(), count+1, (page-1)*count)
	if err != nil {
		return nil, -1, err
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
			return nil, -1, err
		}
		blog := model.NewBlog(blogID, blogWriter, blogTitle, blogText, blogCreatedAt, blogUpdatedAt)
		blogList = append(blogList, blog)
	}

	if err := rows.Err(); err != nil {
		return nil, -1, err
	}

	nextPage := -1
	if len(blogList) > count {
		blogList = blogList[:count]
		nextPage = page + 1
	}

	return blogList, nextPage, nil
}

func (s *BlogStorage) FindByID(id uint64) (*model.Blog, error) {
	stmt := s.db.MustPrepare(`SELECT id, user, title, text, created_at, updated_at FROM blog WHERE id = ?`)
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

func (s *BlogStorage) Update(blog *model.Blog) error {
	stmt := s.db.MustPrepare(`UPDATE blog SET title = ?, text = ? WHERE id = ? and user = ?`)
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
		return storage.ErrNoChange
	}
	return nil
}

func (s *BlogStorage) SetBlogCategory(blog *model.Blog, category *model.Category) error {
	stmt := s.db.MustPrepare(`
		INSERT INTO blog_category (blog, category)
		VALUES(?, ?)
		ON DUPLICATE KEY UPDATE category = ?
	`)
	defer stmt.Close()

	_, err := stmt.Exec(blog.ID(), category.ID(), category.ID())
	return err
}

func (s *BlogStorage) RemoveBlogCategory(blog *model.Blog) error {
	stmt := s.db.MustPrepare(`DELETE FROM blog_category WHERE blog = ?`)
	defer stmt.Close()

	res, err := stmt.Exec(blog.ID())
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrNoRows
	}

	return nil
}
