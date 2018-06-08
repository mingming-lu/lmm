package model

import (
	"lmm/api/testing"
	"time"
)

func TestNewBlog(t *testing.T) {
	tester := testing.NewTester(t)

	created := time.Now()
	updated := time.Now()
	blog := NewBlog(uint64(100), uint64(101), "title", "text", created, updated)

	tester.Is(uint64(100), blog.ID())
	tester.Is(uint64(101), blog.UserID())
	tester.Is("title", blog.Title())
	tester.Is("text", blog.Text())
	tester.Is(created, blog.CreatedAt())
	tester.Is(updated, blog.UpdatedAt())
}

func TestBlogUpdateTitle(t *testing.T) {
	tester := testing.NewTester(t)

	now := time.Now()
	blog := NewBlog(uint64(1), uint64(1), "title", "text", now, now)
	blog.UpdateTitle("new title")

	tester.Is("new title", blog.Title())
	tester.Is("text", blog.Text())
	tester.True(blog.UpdatedAt().After(blog.CreatedAt()))
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
