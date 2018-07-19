package model

import (
	"lmm/api/context/blog/domain"
	"lmm/api/testing"
)

func TestNewTag_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	tag, err := NewTag(123, 234, "tag name")
	t.NoError(err)
	t.Is(uint64(123), tag.ID())
	t.Is(uint64(234), tag.BlogID())
	t.Is("tag name", tag.Name())
}

func TestNewTag_InvalidName(tt *testing.T) {
	t := testing.NewTester(tt)
	tag, err := NewTag(123, 234, "@#?")
	t.IsError(domain.ErrInvalidTagName, err)
	t.Nil(tag)
}

func TestUpdateTagName_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	tag, err := NewTag(123, 234, "tag name")
	t.NoError(err)
	t.NoError(tag.UpdateName("new tag name"))
}

func TestUpdateTagName_InvalidName(tt *testing.T) {
	t := testing.NewTester(tt)
	tag, err := NewTag(333, 111, "tag name")
	t.NoError(err)
	t.IsError(domain.ErrInvalidTagName, tag.UpdateName(""))
}
