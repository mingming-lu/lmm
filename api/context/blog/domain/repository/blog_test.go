package repository

import (
	accountFactory "lmm/api/context/account/domain/factory"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestAddBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)
	err := repo.Add(blog)

	t.NoError(err)

	var (
		userID    uint64
		blogTitle string
		blogText  string
	)
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT user, title, text FROM blog WHERE id = ?`)
	defer stmt.Close()

	stmt.QueryRow(blog.ID()).Scan(&userID, &blogTitle, &blogText)

	t.Is(user.ID(), userID)
	t.Is(blog.Title(), blogTitle)
	t.Is(blog.Text(), blogText)
}

func TestAddBlog_DuplicateTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)

	err := repo.Add(blog)
	t.NoError(err)

	blogWithSameTitle, _ := factory.NewBlog(user.ID(), title, text)
	err = repo.Add(blogWithSameTitle)
	t.Regexp(`Duplicate entry '[\w\d-]+' for key 'title'`, err.Error())
}

func TestFindBlogByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)
	err := repo.Add(blog)

	t.NoError(err)

	blogFound, err := repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(blog.ID(), blogFound.ID())
	t.Is(blog.Title(), blogFound.Title())
	t.Is(blog.Text(), blogFound.Text())
}

func TestFindBlogByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	blog, err := repo.FindByID(uint64(777))
	t.Error(err, "sql: no rows in result set")
	t.Nil(blog)
}
