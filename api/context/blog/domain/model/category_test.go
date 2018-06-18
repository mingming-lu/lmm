package model

import (
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestNewCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := uuid.New()[:31]
	category, err := NewCategory(id, name)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(name, category.Name())
}

func TestNewCategory_ChinesesCharacters(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := "这里全部都是中文汉字哦"
	category, err := NewCategory(id, name)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(name, category.Name())
}

func TestNewCategory_JapaneseCharacters(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := "覚えていますか"
	category, err := NewCategory(id, name)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(name, category.Name())
}

func TestNewCategory_EmptyName(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := ""
	category, err := NewCategory(id, name)

	t.Is(ErrInvalidCategoryName, err)
	t.Nil(category)
}

func TestNewCategory_NameTooLong(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := uuid.New() + uuid.New()
	category, err := NewCategory(id, name)

	t.Is(ErrInvalidCategoryName, err)
	t.Nil(category)
}

func TestUpdateCategoryName_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := uuid.New()[:31]
	category, err := NewCategory(id, name)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(name, category.Name())

	newName := uuid.New()[:31]
	t.NoError(category.UpdateName(newName))
	t.Is(id, category.ID())
	t.Not(name, category.Name())
	t.Is(newName, category.Name())
}

func TestUpdateCategoryName_Failure(tt *testing.T) {
	t := testing.NewTester(tt)

	id := testing.GenerateID()
	name := uuid.New()[:31]
	category, err := NewCategory(id, name)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(name, category.Name())

	newName := ""
	t.Is(ErrInvalidCategoryName, category.UpdateName(newName))
	t.Is(id, category.ID())
	t.Is(name, category.Name())
}
