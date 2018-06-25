package appservice

import (
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestGetBlogListByPage_DefaultCount(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)

	for i := 0; i < 12; i++ {
		_, err := app.blogService.PostBlog(user.ID(), uuid.New(), uuid.New())
		if err != nil {
			panic(err)
		}
	}

	blogListPage, err := app.GetBlogListByPage("", "")

	t.NoError(err)
	t.Is(2, blogListPage.NextPage)
	t.Is(10, len(blogListPage.Blog))

	blogListPage, err = app.GetBlogListByPage("", "2")
	t.NoError(err)
	t.Is(-1, blogListPage.NextPage)
	t.Is(2, len(blogListPage.Blog))

	blogListPage, err = app.GetBlogListByPage("", "3")
	t.NoError(err)
	t.Is(-1, blogListPage.NextPage)
	t.Is(0, len(blogListPage.Blog))
}

// func TestFindAllBlog_GivenCount(tt *testing.T) {
// 	testing.Lock()
// 	defer testing.Unlock()
//
// 	testing.InitTable("blog")
//
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	blogIDs := make([]uint64, 0)
// 	for i := 0; i < 5; i++ {
// 		blogID, err := app.PostNewBlog(user.ID(), uuid.New(), uuid.New())
// 		if err != nil {
// 			panic(err)
// 		}
// 		blogIDs = append(blogIDs, blogID)
// 	}
//
// 	blogList, hasNextPage, err := app.FindAllBlog("3", "")
// 	t.NoError(err)
// 	t.True(hasNextPage)
// 	t.Is(3, len(blogList))
//
// 	blogList, hasNextPage, err = app.FindAllBlog("3", "2")
// 	t.NoError(err)
// 	t.False(hasNextPage)
// 	t.Is(2, len(blogList))
// }
//
// func TestFindBlogByID_Success(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	title, text := uuid.New(), uuid.New()
// 	blogID, err := app.PostNewBlog(user.ID(), title, text)
//
// 	t.NoError(err)
//
// 	blog, err := app.FindBlogByID(fmt.Sprintf("%d", blogID))
// 	t.NoError(err)
// 	t.Is(blogID, blog.ID())
// 	t.Is(title, blog.Title())
// 	t.Is(text, blog.Text())
// }
//
// func TestFindBlogByID_InvalidID(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	blog, err := app.FindBlogByID("NAN")
// 	t.Is(ErrNoSuchBlog, err)
// 	t.Nil(blog)
// }
// func TestFindBlogByID_NotFound(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	blog, err := app.FindBlogByID("112233")
// 	t.Is(ErrNoSuchBlog, err)
// 	t.Nil(blog)
// }
//
// func TestEditBlog_Success(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	title, text := uuid.New(), uuid.New()
// 	blog, _ := factory.NewBlog(user.ID(), title, text)
// 	repo.Add(blog)
//
// 	blog, err := repo.FindByID(blog.ID())
// 	t.NoError(err)
// 	t.Is(title, blog.Title())
// 	t.Is(text, blog.Text())
//
// 	newTitle, newText := uuid.New(), uuid.New()
//
// 	app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), newTitle, newText)
//
// 	blog, err = repo.FindByID(blog.ID())
// 	t.NoError(err)
// 	t.Is(newTitle, blog.Title())
// 	t.Is(newText, blog.Text())
// }
//
// func TestEditBlog_NoPermission(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	title, text := uuid.New(), uuid.New()
// 	blog, _ := factory.NewBlog(user.ID(), title, text)
// 	repo.Add(blog)
//
// 	suspicious, _ := accountFactory.NewUser(uuid.New()[:31], uuid.New())
// 	accountRepository.New(testing.DB()).Add(suspicious)
//
// 	err := app.EditBlog(suspicious.ID(), fmt.Sprintf("%d", blog.ID()), "new title", "new text")
// 	t.Is(ErrNoPermission, err)
// }
//
// func TestEditBlog_NoSuchBlog(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	title, text := uuid.New(), uuid.New()
// 	blog, _ := factory.NewBlog(user.ID(), title, text)
//
// 	// notice that I didn' save that blog and here I reverse the title and the text to exclude ErrBlogNoChange
// 	err := app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), text, title)
// 	t.Is(ErrNoSuchBlog, err)
// }
//
// func TestEditBlog_NoChange(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	title, text := uuid.New(), uuid.New()
// 	blog, _ := factory.NewBlog(user.ID(), title, text)
// 	repo.Add(blog)
//
// 	err := app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), blog.Title(), blog.Text())
// 	t.Is(ErrBlogNoChange, err)
// }
//
// func TestEditBlog_EmptyTitle(tt *testing.T) {
// 	t := testing.NewTester(tt)
// 	repo := repository.NewBlogRepository(testing.DB())
// 	app := NewBlogApp(repo)
//
// 	title, text := uuid.New(), uuid.New()
// 	blog, _ := factory.NewBlog(user.ID(), title, text)
// 	repo.Add(blog)
//
// 	err := app.EditBlog(user.ID(), fmt.Sprintf("%d", blog.ID()), "", blog.Text())
// 	t.Is(ErrEmptyBlogTitle, err)
// }
