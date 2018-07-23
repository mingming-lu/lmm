package infra

import (
	accountFactory "lmm/api/context/account/domain/factory"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/storage"
	"lmm/api/testing"
	"lmm/api/utils/strings"
	"lmm/api/utils/uuid"
	"math/rand"
	"time"
)

func TestAddTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	tag, err := factory.NewTag(12, uuid.New()[:31])
	t.NoError(err)
	t.NoError(repo.Add(tag))

	tagFound, err := repo.selectRow(
		`SELECT id, blog, name FROM tag WHERE id = ? AND blog = ? AND name = ?`,
		tag.ID(), tag.BlogID(), tag.Name(),
	)
	t.NoError(err)
	t.Is(tag, tagFound)
}

func TestAddTag_DuplicateBlogTag(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("tag")

	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	name := uuid.New()[:31]
	if tag, err := factory.NewTag(12, name); true {
		t.NoError(err)
		t.NoError(repo.Add(tag))
	}

	tag, err := factory.NewTag(12, name)
	t.NoError(err)

	err = repo.Add(tag)
	t.IsError(domain.ErrDuplicateTagName, err)
}

func TestFindTagByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	name := uuid.New()[:31]
	tag, err := factory.NewTag(11, name)
	t.NoError(err)
	t.NoError(repo.Add(tag))

	tagFound, err := repo.FindByID(tag.ID())
	t.NoError(err)
	t.Is(tag, tagFound)
}

func TestFindTagByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	name := uuid.New()[:31]
	tag, err := factory.NewTag(11, name)
	t.NoError(err)

	tagFound, err := repo.FindByID(tag.ID())
	t.IsError(storage.ErrNoRows, err)
	t.Nil(tagFound)
}

func TestFindAllTags_Empty(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("tag")
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	tags, err := repo.FindAll()
	t.NoError(err)
	t.NotNil(tags)
	t.Is(0, len(tags))
}

func TestFindAllTags_SortByName(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("tag")

	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	for _, name := range "gafcedbh" {
		blogID := uint64(rand.Int63n(10) + 1)
		tag, err := factory.NewTag(blogID, string(name))
		t.NoError(err)
		t.NoError(repo.Add(tag))
	}

	tags, err := repo.FindAll()
	t.NoError(err)

	names := make([]string, len(tags))
	for index, tag := range tags {
		names[index] = tag.Name()
	}
	t.Is("abcdefgh", strings.Join("", names))
}

func TestFindAllTagsByBlog_Empty(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	var blog *model.Blog
	// preparing blog data
	if user, err := accountFactory.NewUser(uuid.New()[:31], uuid.New()); true {
		t.NoError(err)
		title, text := uuid.New(), uuid.New()
		blog, err = factory.NewBlog(user.ID(), title, text)
		t.NoError(err)
		blogRepo := NewBlogStorage(testing.DB())
		t.NoError(blogRepo.Add(blog))
	}

	tags, err := repo.FindAllByBlog(blog)
	t.NoError(err)
	t.NotNil(tags)
	t.Is(0, len(tags))
}

func TestFindAllTagsByBlog_SortByName(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	var blog *model.Blog
	// preparing blog data
	if user, err := accountFactory.NewUser(uuid.New()[:31], uuid.New()); true {
		t.NoError(err)
		title, text := uuid.New(), uuid.New()
		blog, err = factory.NewBlog(user.ID(), title, text)
		t.NoError(err)
		blogRepo := NewBlogStorage(testing.DB())
		t.NoError(blogRepo.Add(blog))
	}

	// preparing tag data
	for _, name := range "liknjm" {
		tag, err := factory.NewTag(blog.ID(), string(name))
		t.NoError(err)
		t.NoError(repo.Add(tag))
	}

	tags, err := repo.FindAllByBlog(blog)
	t.NoError(err)
	names := make([]string, len(tags))
	for index, tag := range tags {
		names[index] = tag.Name()
	}
	t.Is("ijklmn", strings.Join("", names))
}

func TestFindAllTagsByBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	now := time.Now()
	dummyBlog := model.NewBlog(1233, 123, "dummy", "for testing", now, now)

	tags, err := repo.FindAllByBlog(dummyBlog)
	t.NoError(err)
	t.NotNil(tags)
	t.Is(0, len(tags))
}

func TestUpdateTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	tag, err := factory.NewTag(333, "aaa")
	t.NoError(err)
	t.NoError(repo.Add(tag))

	if tagFound, err := repo.FindByID(tag.ID()); true {
		t.NoError(err)
		t.Is("aaa", tagFound.Name())
	}

	t.NoError(tag.UpdateName("bbb"))
	t.NoError(repo.Update(tag))

	if tagFound, err := repo.FindByID(tag.ID()); true {
		t.NoError(err)
		t.Is("bbb", tagFound.Name())
	}
}

func TestUpdateTag_NoSuchTag(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	tag, err := factory.NewTag(333, "abc")
	t.NoError(err)

	t.NoError(tag.UpdateName("bbb"))
	t.IsErrorMsg("rows affected is expected to be 1 but got 0", repo.Update(tag))
}

func TestRemoveTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	tag, err := factory.NewTag(333, "aaa")
	t.NoError(err)
	t.NoError(repo.Add(tag))

	t.NoError(repo.Remove(tag))

	tagFound, err := repo.FindByID(tag.ID())
	t.IsError(storage.ErrNoRows, err)
	t.Nil(tagFound)
}

func TestRemoveTag_NoSuchTag(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewTagStorage(testing.DB())

	tag, err := factory.NewTag(333, "foobar")
	t.NoError(err)
	t.NoError(repo.Add(tag))

	t.NoError(repo.Remove(tag))
	t.IsErrorMsg("rows affected is expected to be 1 but got 0", repo.Remove(tag))
}
