package infra

import (
	"fmt"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/storage"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
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
	key, entry, ok := storage.CheckErrorDuplicate(err.Error())
	t.True(ok)
	t.Is(fmt.Sprintf("%d-%s", tag.BlogID(), tag.Name()), entry)
	t.Is("blog_tag", key)
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
