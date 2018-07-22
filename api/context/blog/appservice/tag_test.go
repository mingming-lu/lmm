package appservice

import (
	"fmt"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/infra"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestAddNewTagToBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	blogRepo.Add(blog)
	t.NoError(err)

	tagName := uuid.New()[:31]
	t.NoError(app.AddNewTagToBlog(user, fmt.Sprint(blog.ID()), tagName))

	tags, err := app.tagRepository.FindAllByBlog(blog)
	t.NoError(err)
	t.NotNil(tags)
	t.Is(1, len(tags))
	t.Is(tagName, tags[0].Name())
}

func TestAddNewTagToBlog_InvalidBlogID(tt *testing.T) {
	t := testing.NewTester(tt)
	t.IsError(
		domain.ErrNoSuchBlog,
		app.AddNewTagToBlog(user, "invalid blog id", "valid tag name"),
	)
}

func TestAddNewTagToBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)

	t.IsError(
		domain.ErrNoSuchBlog,
		app.AddNewTagToBlog(user, fmt.Sprint(blog.ID()), "valid tag name"),
	)
}

func TestAddNewTagToBlog_InvalidTagName(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	t.IsError(
		domain.ErrInvalidTagName,
		app.AddNewTagToBlog(user, fmt.Sprint(blog.ID()), uuid.New()),
	)
}

func TestAddNewTagToBlog_Duplicate(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tagName := uuid.New()[:31]
	t.NoError(app.AddNewTagToBlog(user, fmt.Sprint(blog.ID()), tagName))
	t.IsError(
		domain.ErrDuplicateTagName,
		app.AddNewTagToBlog(user, fmt.Sprint(blog.ID()), tagName),
	)
}

func TestUpdateTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tagName := uuid.New()[:31]
	tag, err := factory.NewTag(blog.ID(), tagName)
	t.NoError(err)
	t.NoError(app.tagRepository.Add(tag))

	newTagName := uuid.New()[:31]
	t.NoError(app.UpdateBlogTag(user, fmt.Sprint(tag.ID()), newTagName))
}

func TestUpdateTag_InvalidTagID(tt *testing.T) {
	t := testing.NewTester(tt)
	t.IsError(domain.ErrNoSuchTag, app.UpdateBlogTag(user, "invalid tag id", "dummy name"))
}

func TestUpdataTag_NoSuchTag(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tagName := uuid.New()[:31]
	tag, err := factory.NewTag(blog.ID(), tagName)
	t.NoError(err)
	// t.NoError(app.tagRepository.Add(tag))

	newTagName := uuid.New()[:31]
	t.IsError(domain.ErrNoSuchTag, app.UpdateBlogTag(user, fmt.Sprint(tag.ID()), newTagName))
}

func TestUpdateTag_InvalidTagName(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tagName := uuid.New()[:31]
	tag, err := factory.NewTag(blog.ID(), tagName)
	t.NoError(err)
	t.NoError(app.tagRepository.Add(tag))

	t.IsError(domain.ErrInvalidTagName, app.UpdateBlogTag(user, fmt.Sprint(tag.ID()), uuid.New()))
}

func TestRemoveBlogTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tagName := uuid.New()[:31]
	tag, err := factory.NewTag(blog.ID(), tagName)
	t.NoError(err)
	t.NoError(app.tagRepository.Add(tag))

	t.NoError(app.RemoveBlogTag(user, fmt.Sprint(tag.ID())))
}

func TestRemoveBlogTag_InvalidTagID(tt *testing.T) {
	t := testing.NewTester(tt)
	t.IsError(domain.ErrNoSuchTag, app.UpdateBlogTag(user, "invalid tag id", "dummy name"))
}

func TestRemoveBlogTag_NoSuchTag(tt *testing.T) {
	t := testing.NewTester(tt)
	blogRepo := infra.NewBlogStorage(testing.DB())

	blog, err := factory.NewBlog(user.ID(), uuid.New()[:31], uuid.New())
	t.NoError(err)
	t.NoError(blogRepo.Add(blog))

	tagName := uuid.New()[:31]
	tag, err := factory.NewTag(blog.ID(), tagName)
	t.NoError(err)
	// t.NoError(app.tagRepository.Add(tag))

	t.IsError(domain.ErrNoSuchTag, app.RemoveBlogTag(user, fmt.Sprint(tag.ID())))
}
