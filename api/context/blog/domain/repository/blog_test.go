package repository

import (
	accountFactory "lmm/api/context/account/domain/factory"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestAddBlog(t *testing.T) {
	tester := testing.NewTester(t)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)
	err := repo.Add(blog)

	tester.NoError(err)

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

	tester.Is(user.ID(), userID)
	tester.Is(blog.Title(), blogTitle)
	tester.Is(blog.Text(), blogText)
}
