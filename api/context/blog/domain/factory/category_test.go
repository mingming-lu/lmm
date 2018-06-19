package factory

import (
	"lmm/api/context/blog/domain/model"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestNewCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	name := uuid.New()[:31]
	category, err := NewCategory(name)

	t.NoError(err)
	t.True(category.ID() > 0)
	t.Is(name, category.Name())
}

func TestNewCategory_Failure(tt *testing.T) {
	t := testing.NewTester(tt)

	name := ""
	category, err := NewCategory(name)

	t.Is(model.ErrInvalidCategoryName, err)
	t.Nil(category)
}
