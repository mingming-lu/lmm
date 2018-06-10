package appservice

import (
	"fmt"
	accountFactory "lmm/api/context/account/domain/factory"
	account "lmm/api/context/account/domain/model"
	accountRepository "lmm/api/context/account/domain/repository"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"os"
)

var user *account.User

func TestMain(m *testing.M) {
	name, password := uuid.New()[:31], uuid.New()
	user, _ = accountFactory.NewUser(name, password)
	accountRepository.New().Add(user)

	code := m.Run()
	os.Exit(code)
}

func TestPostNewBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := New(repo)

	title, text := uuid.New(), uuid.New()
	blogID, err := app.PostNewBlog(user.ID(), title, text)

	t.NoError(err)
	t.True(blogID > uint64(0))

	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT title, text FROM blog WHERE id = ?`)
	defer stmt.Close()

	var (
		blogTitle string
		blogText  string
	)
	err = stmt.QueryRow(blogID).Scan(&blogTitle, &blogText)
	t.NoError(err)
	t.Is(title, blogTitle)
	t.Is(text, blogText)
}

func TestPostNewBlog_DuplicateTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := New(repo)

	title, text := uuid.New(), uuid.New()
	_, err := app.PostNewBlog(user.ID(), title, text)
	t.NoError(err)

	_, err = app.PostNewBlog(user.ID(), title, text)
	t.Is(ErrBlogTitleDuplicated, err)
}

func TestFindBlogByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := New(repo)

	title, text := uuid.New(), uuid.New()
	blogID, err := app.PostNewBlog(user.ID(), title, text)

	t.NoError(err)

	blog, err := app.FindBlogByID(fmt.Sprintf("%d", blogID))
	t.NoError(err)
	t.Is(blogID, blog.ID())
	t.Is(title, blog.Title())
	t.Is(text, blog.Text())
}

func TestFindBlogByID_InvalidID(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := New(repo)

	blog, err := app.FindBlogByID("NAN")
	t.Is(ErrNoSuchBlog, err)
	t.Nil(blog)
}
func TestFindBlogByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := New(repo)

	blog, err := app.FindBlogByID("112233")
	t.Is(ErrNoSuchBlog, err)
	t.Nil(blog)
}
