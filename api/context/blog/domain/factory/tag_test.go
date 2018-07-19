package factory

import (
	"lmm/api/context/blog/domain"
	"lmm/api/testing"
)

func TestNewTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	tag, err := NewTag(323, "tag")
	t.NoError(err)
	t.Is(uint64(323), tag.BlogID())
	t.Is("tag", tag.Name())
}

func TestNewTag_InvalidName(tt *testing.T) {
	t := testing.NewTester(tt)
	tag, err := NewTag(323, "@@")
	t.IsError(domain.ErrInvalidTagName, err)
	t.Nil(tag)
}
