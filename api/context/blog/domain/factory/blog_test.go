package factory

import "lmm/api/testing"

func TestNewBlog(tt *testing.T) {
	t := testing.NewTester(tt)

	blog, err := NewBlog(uint64(1), "title", "text")

	t.NoError(err)
	t.True(blog.ID() > uint64(0))
	t.Is(uint64(1), blog.UserID())
	t.Is("title", blog.Title())
	t.Is("text", blog.Text())
	t.Is(blog.CreatedAt(), blog.UpdatedAt())
}
