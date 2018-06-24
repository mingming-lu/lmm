package repository

import (
	accountFactory "lmm/api/context/account/domain/factory"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/repository"
	"lmm/api/storage"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"sort"
)

func TestNewAddCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()

	err := repo.Add(category)

	t.NoError(err)
}

func TestNewAddCategory_Duplicated(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()

	err := repo.Add(category)

	t.NoError(err)

	err = repo.Add(category)
	t.True(repository.ErrPatternDuplicate.Match([]byte(err.Error())))
}

func TestUpdateCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()

	err := repo.Add(category)
	t.NoError(err)
}

func TestUpdateCategory_NoSuchCategory(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()

	err := repo.Update(category)
	t.Isa(storage.ErrNoChange, err)
}

func TestFindAllCategories_Success(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()
	testing.InitTable("category")

	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	for i := 0; i < 10; i++ {
		repo.Add(newCategory())
	}

	categories, err := repo.FindAll()

	t.NoError(err)
	t.Is(10, len(categories))

	names := make([]string, len(categories))
	for index, category := range categories {
		names[index] = category.Name()
	}

	t.True(sort.StringsAreSorted(names))
}

func TestFindAllCategories_NoOne(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("category")

	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	categories, err := repo.FindAll()
	t.NoError(err)
	t.Is(0, len(categories))
}

func TestFindCategoryByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()
	repo.Add(category)

	categoryFound, err := repo.FindByID(category.ID())

	t.NoError(err)
	t.Is(category.ID(), categoryFound.ID())
	t.Is(category.Name(), categoryFound.Name())
}

func TestFindCategoryByID_NoSuch(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()
	// repo.Add(category)

	categoryFound, err := repo.FindByID(category.ID())

	t.Is(storage.ErrNoRows, err)
	t.Nil(categoryFound)
}

func TestFindCategoryByBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	insertBlogCategory := testing.DB().MustPrepare(`INSERT INTO blog_category (blog, category) VALUES(?, ?)`)
	defer insertBlogCategory.Close()

	category := newCategory()
	repo.Add(category)

	user, _ := accountFactory.NewUser("username", "userpassword")
	blog, _ := factory.NewBlog(user.ID(), "blog title", "blog text")

	insertBlogCategory.Exec(blog.ID(), category.ID())

	categoryFound, err := repo.FindByBlog(blog)
	t.NoError(err)
	t.Is(category.ID(), categoryFound.ID())
}

func TestFindCategoryByBlog_CategoryNotSet(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()
	repo.Add(category)

	user, _ := accountFactory.NewUser("username", "userpassword")
	blog, _ := factory.NewBlog(user.ID(), "blog title", "blog text")

	categoryFound, err := repo.FindByBlog(blog)
	t.Is(storage.ErrNoRows, err)
	t.Nil(categoryFound)
}

func TestFindCategoryByBlog_NoSuchCategory(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	insertBlogCategory := testing.DB().MustPrepare(`INSERT INTO blog_category (blog, category) VALUES(?, ?)`)
	defer insertBlogCategory.Close()

	category := newCategory()

	user, _ := accountFactory.NewUser("username", "userpassword")
	blog, _ := factory.NewBlog(user.ID(), "blog title", "blog text")

	insertBlogCategory.Exec(blog.ID(), category.ID())

	categoryFound, err := repo.FindByBlog(blog)
	t.Is(storage.ErrNoRows, err)
	t.Nil(categoryFound)
}

func TestRemoveCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()
	repo.Add(category)

	err := repo.Remove(category)

	t.NoError(err)

	category, err = repo.FindByID(category.ID())
	t.Is(storage.ErrNoRows, err)
	t.Nil(category)
}

func TestRemoveCategory_NoSuch(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()
	// repo.Add(category)

	err := repo.Remove(category)

	t.Is(storage.ErrNoRows, err)
}

func TestRemoveCategory_Removed(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewCategoryRepository(testing.DB())

	category := newCategory()
	repo.Add(category)

	err := repo.Remove(category)

	t.NoError(err)

	t.Is(storage.ErrNoRows, repo.Remove(category))
}

func newCategory() *model.Category {
	name := uuid.New()[:31]
	category, err := factory.NewCategory(name)

	if err != nil {
		panic(err)
	}

	return category
}
