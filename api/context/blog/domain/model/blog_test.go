package model

import (
	"lmm/api/testing"
	"time"
)

func TestNewBlog(tt *testing.T) {
	t := testing.NewTester(tt)

	created := time.Now()
	updated := time.Now()
	blog := NewBlog(uint64(100), uint64(101), "title", "text", created, updated)

	t.Is(uint64(100), blog.ID())
	t.Is(uint64(101), blog.UserID())
	t.Is("title", blog.Title())
	t.Is("text", blog.Text())
	t.Is(created, blog.CreatedAt())
	t.Is(updated, blog.UpdatedAt())
}

func TestBlogUpdateTitle(tt *testing.T) {
	t := testing.NewTester(tt)

	now := time.Now()
	blog := NewBlog(uint64(1), uint64(1), "title", "text", now, now)
	blog.UpdateTitle("new title")

	t.Is("new title", blog.Title())
	t.Is("text", blog.Text())
	t.True(blog.UpdatedAt().After(blog.CreatedAt()))
}

func TestBlogUpdateTitle_NoChange(tt *testing.T) {
	t := testing.NewTester(tt)

	now := time.Now()
	blog := NewBlog(uint64(1), uint64(1), "title", "text", now, now)
	blog.UpdateTitle("title")

	t.Is("title", blog.Title())
	t.True(blog.UpdatedAt().Equal(blog.CreatedAt()))
}

func TestBlogUpdateText(t *testing.T) {
	tester := testing.NewTester(t)

	now := time.Now()
	blog := NewBlog(uint64(1), uint64(1), "title", "text", now, now)
	blog.UpdateText("new text")

	tester.Is("new text", blog.Text())
	tester.Is("title", blog.Title())
	tester.True(blog.UpdatedAt().After(blog.CreatedAt()))
}

func TestBlogUpdateText_NoChange(tt *testing.T) {
	t := testing.NewTester(tt)

	now := time.Now()
	blog := NewBlog(uint64(1), uint64(1), "title", "text", now, now)
	blog.UpdateText("text")

	t.Is("text", blog.Text())
	t.True(blog.UpdatedAt().Equal(blog.CreatedAt()))
}
