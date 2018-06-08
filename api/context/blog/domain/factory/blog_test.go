package factory

import "lmm/api/testing"

func TestNewBlog(t *testing.T) {
	tester := testing.NewTester(t)

	blog, err := NewBlog(uint64(1), "title", "text")

	tester.NoError(err)
	tester.True(blog.ID() > uint64(0))
	tester.Is(uint64(1), blog.UserID())
	tester.Is("title", blog.Title())
	tester.Is("text", blog.Text())
	tester.Is(blog.CreatedAt(), blog.UpdatedAt())
}
