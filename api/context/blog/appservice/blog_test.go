package appservice

import (
	"fmt"
	accountFactory "lmm/api/context/account/domain/factory"
	accountRepository "lmm/api/context/account/domain/repository"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"time"
)

func TestPostNewBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

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
	app := NewBlogApp(repo)

	title, text := uuid.New(), uuid.New()
	_, err := app.PostNewBlog(user.ID(), title, text)
	t.NoError(err)

	_, err = app.PostNewBlog(user.ID(), title, text)
	t.Is(ErrBlogTitleDuplicated, err)
}

func TestPostNewBlog_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	blogID, err := app.PostNewBlog(user.ID(), "", uuid.New())
	t.Is(ErrEmptyBlogTitle, err)
	t.Is(uint64(0), blogID)
}

func TestFindAllBlog_OrderByCreatedTime(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	app := NewBlogApp(repository.NewBlogRepository())

	app.PostNewBlog(user.ID(), uuid.New(), uuid.New())
	time.Sleep(1 * time.Second)
	app.PostNewBlog(user.ID(), uuid.New(), uuid.New())

	blogList, page, hasNextPage, err := app.FindAllBlog("", "")
	t.NoError(err)
	t.False(hasNextPage)
	t.Is(1, page)
	t.Is(2, len(blogList))
	t.True(blogList[0].CreatedAt().After(blogList[1].CreatedAt()))
}

func TestFindAllBlog_DefaultCount(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	for i := 0; i < 12; i++ {
		_, err := app.PostNewBlog(user.ID(), uuid.New(), uuid.New())
		if err != nil {
			panic(err)
		}
	}

	blogList, page, hasNextPage, err := app.FindAllBlog("", "")

	t.NoError(err)
	t.True(hasNextPage)
	t.Is(1, page)
	t.Is(10, len(blogList))

	blogList, page, hasNextPage, err = app.FindAllBlog("", fmt.Sprintf("%d", page+1))
	t.NoError(err)
	t.False(hasNextPage)
	t.Is(2, page)
	t.Is(2, len(blogList))

	blogList, page, hasNextPage, err = app.FindAllBlog("", fmt.Sprintf("%d", page+1))
	t.NoError(err)
	t.False(hasNextPage)
	t.Is(3, page)
	t.Is(0, len(blogList))
}

func TestFindAllBlog_GivenCount(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	blogIDs := make([]uint64, 0)
	for i := 0; i < 5; i++ {
		blogID, err := app.PostNewBlog(user.ID(), uuid.New(), uuid.New())
		if err != nil {
			panic(err)
		}
		blogIDs = append(blogIDs, blogID)
	}

	blogList, page, hasNextPage, err := app.FindAllBlog("3", "")
	t.NoError(err)
	t.True(hasNextPage)
	t.Is(1, page)
	t.Is(3, len(blogList))

	blogList, page, hasNextPage, err = app.FindAllBlog("3", fmt.Sprintf("%d", page+1))
	t.NoError(err)
	t.False(hasNextPage)
	t.Is(2, page)
	t.Is(2, len(blogList))
}

func TestFindBlogByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

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
	app := NewBlogApp(repo)

	blog, err := app.FindBlogByID("NAN")
	t.Is(ErrNoSuchBlog, err)
	t.Nil(blog)
}
func TestFindBlogByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	blog, err := app.FindBlogByID("112233")
	t.Is(ErrNoSuchBlog, err)
	t.Nil(blog)
}

func TestEditBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	blog, err := repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(title, blog.Title())
	t.Is(text, blog.Text())

	newTitle, newText := uuid.New(), uuid.New()

	app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), newTitle, newText)

	blog, err = repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(newTitle, blog.Title())
	t.Is(newText, blog.Text())
}

func TestEditBlog_NoPermission(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	suspicious, _ := accountFactory.NewUser(uuid.New()[:31], uuid.New())
	accountRepository.New().Add(suspicious)

	err := app.EditBlog(suspicious.ID(), fmt.Sprintf("%d", blog.ID()), "new title", "new text")
	t.Is(ErrNoPermission, err)
}

func TestEditBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)

	// notice that I didn' save that blog and here I reverse the title and the text to exclude ErrBlogNoChange
	err := app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), text, title)
	t.Is(ErrNoSuchBlog, err)
}

func TestEditBlog_NoChange(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	err := app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), blog.Title(), blog.Text())
	t.Is(ErrBlogNoChange, err)
}

func TestEditBlog_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository()
	app := NewBlogApp(repo)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	err := app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), "", blog.Text())
	t.Is(ErrEmptyBlogTitle, err)
}
