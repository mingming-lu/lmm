package appservice

import (
	"fmt"
	accountFactory "lmm/api/context/account/domain/factory"
	accountRepository "lmm/api/context/account/domain/repository"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/service"
	"lmm/api/context/blog/repository"
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

func TestFindAllBlog_GivenCount(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)

	for i := 0; i < 5; i++ {
		_, err := app.blogService.PostBlog(user.ID(), uuid.New(), uuid.New())
		if err != nil {
			panic(err)
		}
	}

	blogListPage, err := app.GetBlogListByPage("3", "")
	t.NoError(err)
	t.Is(2, blogListPage.NextPage)
	t.Is(3, len(blogListPage.Blog))

	blogListPage, err = app.GetBlogListByPage("3", "2")
	t.NoError(err)
	t.Is(-1, blogListPage.NextPage)
	t.Is(2, len(blogListPage.Blog))
}

func TestFindBlogByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	blogContent := BlogContent{
		Title: uuid.New(),
		Text:  uuid.New(),
	}
	blogID, err := app.PostNewBlog(user, testing.StructToRequestBody(blogContent))

	t.NoError(err)

	blog, err := app.GetBlogByID(fmt.Sprintf("%d", blogID))
	t.NoError(err)
	t.Is(blogID, blog.ID())
	t.Is(blogContent.Title, blog.Title())
	t.Is(blogContent.Text, blog.Text())
}

func TestFindBlogByID_InvalidID(tt *testing.T) {
	t := testing.NewTester(tt)

	blog, err := app.GetBlogByID("NAN")
	t.Is(service.ErrInvalidBlogID, err)
	t.Nil(blog)
}

func TestFindBlogByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)

	blog, err := app.GetBlogByID("112233")
	t.Is(service.ErrNoSuchBlog, err)
	t.Nil(blog)
}

func TestEditBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	blog, err := repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(title, blog.Title())
	t.Is(text, blog.Text())

	blogContent := BlogContent{
		Title: uuid.New(),
		Text:  uuid.New(),
	}

	t.NoError(app.EditBlog(user, fmt.Sprintf("%d", blog.ID()), testing.StructToRequestBody(blogContent)))

	blog, err = repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(blogContent.Title, blog.Title())
	t.Is(blogContent.Text, blog.Text())
}

func TestEditBlog_NoPermission(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	suspicious, _ := accountFactory.NewUser(uuid.New()[:31], uuid.New())
	accountRepository.New(testing.DB()).Add(suspicious)

	blogContent := BlogContent{
		Title: uuid.New(),
		Text:  uuid.New(),
	}

	t.Is(
		service.ErrNoPermission,
		app.EditBlog(suspicious, fmt.Sprintf("%d", blog.ID()), testing.StructToRequestBody(blogContent)),
	)
}

func TestEditBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)

	blogContent := BlogContent{
		Title: uuid.New(),
		Text:  uuid.New(),
	}

	// notice that I didn' save that blog and here I reverse the title and the text to exclude ErrBlogNoChange
	t.Is(
		service.ErrNoSuchBlog,
		app.EditBlog(user, fmt.Sprintf("%d", blog.ID()), testing.StructToRequestBody(blogContent)),
	)
}

func TestEditBlog_NoChange(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	blogContent := BlogContent{
		Title: blog.Title(),
		Text:  blog.Text(),
	}

	t.Is(
		service.ErrBlogNoChange,
		app.EditBlog(user, fmt.Sprintf("%d", blog.ID()), testing.StructToRequestBody(blogContent)),
	)
}

func TestEditBlog_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewBlogRepository(testing.DB())

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	blogContent := BlogContent{
		Title: "",
		Text:  blog.Text(),
	}

	t.Is(
		service.ErrEmptyBlogTitle,
		app.EditBlog(user, fmt.Sprintf("%d", blog.ID()), testing.StructToRequestBody(blogContent)),
	)
}
