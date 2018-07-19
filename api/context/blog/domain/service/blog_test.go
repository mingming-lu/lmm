package service

import (
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"time"
)

func TestPostBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewBlogStorage(testing.DB())
	service := NewBlogService(repo)

	title, text := uuid.New(), uuid.New()
	blog, err := service.PostBlog(user.ID(), title, text)

	t.NoError(err)
	t.True(blog.ID() > uint64(0))

	stmt := testing.DB().MustPrepare(`SELECT title, text FROM blog WHERE id = ?`)
	defer stmt.Close()

	var (
		blogTitle string
		blogText  string
	)
	err = stmt.QueryRow(blog.ID()).Scan(&blogTitle, &blogText)
	t.NoError(err)
	t.Is(title, blogTitle)
	t.Is(text, blogText)
}

func TestPostBlog_DuplicateTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewBlogStorage(testing.DB())
	service := NewBlogService(repo)

	title, text := uuid.New(), uuid.New()
	blog, err := service.PostBlog(user.ID(), title, text)
	t.NoError(err)
	t.NotNil(blog)
	t.Isa(&model.Blog{}, blog)

	blog, err = service.PostBlog(user.ID(), title, text)
	t.Is(domain.ErrBlogTitleDuplicated, err)
	t.Nil(blog)
}

func TestPostBlog_EmptyTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewBlogStorage(testing.DB())
	service := NewBlogService(repo)

	blog, err := service.PostBlog(user.ID(), "", uuid.New())
	t.Is(domain.ErrEmptyBlogTitle, err)
	t.Nil(blog)
}

func TestFindAllBlog_OrderByCreatedTime(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	repo := infra.NewBlogStorage(testing.DB())
	service := NewBlogService(repo)

	service.PostBlog(user.ID(), uuid.New(), uuid.New())
	time.Sleep(1 * time.Second)
	service.PostBlog(user.ID(), uuid.New(), uuid.New())

	blogList, nextPage, err := service.GetBlogListByPage(10, 1)
	t.NoError(err)
	t.Is(-1, nextPage)
	t.Is(2, len(blogList))
	t.True(blogList[0].CreatedAt().After(blogList[1].CreatedAt()))
}
