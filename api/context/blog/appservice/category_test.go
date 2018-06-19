package appservice

import (
	"fmt"
	"lmm/api/context/blog/domain/repository"
	"lmm/api/domain/factory"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestAddNewCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	name := uuid.New()[:31]
	id, err := app.AddNewCategory(name)

	category, err := repo.FindByID(id)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(name, category.Name())
}

func TestAddNewCategory_InvalidName(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	id, err := app.AddNewCategory("")
	t.Is(ErrInvalidCategoryName, err)
	t.Is(uint64(0), id)
}

func TestAddNewCategory_DuplicatedName(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	name := uuid.New()[:31]
	_, err := app.AddNewCategory(name)
	t.NoError(err)

	_, err = app.AddNewCategory(name)
	t.Is(ErrDuplicateCategoryName, err)
}

func TestUpdateCategoryName_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	name := uuid.New()[:31]
	id, err := app.AddNewCategory(name)
	t.NoError(err)

	newName := uuid.New()[:31]
	t.NoError(app.UpdateCategoryName(fmt.Sprint(id), newName))
}

func TestUpdateCategoryName_InvalidID(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	err := app.UpdateCategoryName("?", "new name")
	t.Is(ErrNoSuchCategory, err)
}

func TestUpdateCategoryName_NoSuchCategory(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	id, _ := factory.Default().GenerateID()
	err := app.UpdateCategoryName(fmt.Sprint(id), "new name")
	t.Is(ErrNoSuchCategory, err)
}

func TestUpdateCategoryName_InvalidCategoryName(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	name := uuid.New()[:31]
	id, err := app.AddNewCategory(name)
	t.NoError(err)

	newName := "#$%"
	err = app.UpdateCategoryName(fmt.Sprint(id), newName)
	t.Is(ErrInvalidCategoryName, err)
}

func TestUpdateCategoryName_CategoryNameNoChanged(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	name := uuid.New()[:31]
	id, err := app.AddNewCategory(name)
	t.NoError(err)

	err = app.UpdateCategoryName(fmt.Sprint(id), name)
	t.Is(ErrCategoryNoChanged, err)
}

func TestRemoveCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	name := uuid.New()[:31]
	id, err := app.AddNewCategory(name)
	t.NoError(err)

	t.NoError(app.Remove(fmt.Sprint(id)))
}

func TestRemoveCategory_InvalidID(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	t.Is(ErrNoSuchCategory, app.Remove("?"))
}

func TestRemoveCategory_NoSuchCategory(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := repository.NewCategoryRepository()
	app := NewCategoryApp(repo)

	id, _ := factory.Default().GenerateID()
	err := app.Remove(fmt.Sprint(id))
	t.Is(ErrNoSuchCategory, err)
}
