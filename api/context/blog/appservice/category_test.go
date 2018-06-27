package appservice

import (
	"fmt"
	"lmm/api/context/blog/domain/model"
	"lmm/api/context/blog/domain/service"
	"lmm/api/context/blog/infra"
	"lmm/api/domain/factory"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestAddNewCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := infra.NewCategoryStorage(testing.DB())

	body := Category{
		Name: uuid.New()[:31],
	}

	id, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))

	category, err := repo.FindByID(id)

	t.NoError(err)
	t.Is(id, category.ID())
	t.Is(body.Name, category.Name())
}

func TestAddNewCategory_EmptyName(tt *testing.T) {
	t := testing.NewTester(tt)

	body := Category{
		Name: "",
	}

	id, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.Is(service.ErrInvalidCategoryName, err)
	t.Is(uint64(0), id)
}

func TestAddNewCategory_DuplicatedName(tt *testing.T) {
	t := testing.NewTester(tt)

	body := Category{
		Name: uuid.New()[:31],
	}

	_, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.NoError(err)

	_, err = app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.Is(service.ErrDuplicateCategoryName, err)
}

func TestUpdateCategoryName_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	body := Category{
		Name: uuid.New()[:31],
	}

	id, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.NoError(err)

	body = Category{
		Name: uuid.New()[:31],
	}

	t.NoError(app.EditCategory(user, fmt.Sprint(id), testing.StructToRequestBody(body)))
}

func TestUpdateCategoryName_InvalidID(tt *testing.T) {
	t := testing.NewTester(tt)

	body := Category{
		Name: uuid.New()[:31],
	}

	_, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.NoError(err)

	body = Category{
		Name: uuid.New()[:31],
	}

	t.Is(
		service.ErrInvalidCategoryID,
		app.EditCategory(user, "invalid id???", testing.StructToRequestBody(body)),
	)
}

func TestUpdateCategoryName_NoSuchCategory(tt *testing.T) {
	t := testing.NewTester(tt)

	id, _ := factory.Default().GenerateID()

	body := Category{
		Name: uuid.New()[:31],
	}

	err := app.EditCategory(user, fmt.Sprint(id), testing.StructToRequestBody(body))

	t.Is(service.ErrNoSuchCategory, err)
}

func TestUpdateCategoryName_InvalidCategoryName(tt *testing.T) {
	t := testing.NewTester(tt)

	body := Category{
		Name: uuid.New()[:31],
	}

	id, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.NoError(err)

	body = Category{
		Name: "%^&oh yeah~!",
	}

	t.Is(
		service.ErrInvalidCategoryName,
		app.EditCategory(user, fmt.Sprint(id), testing.StructToRequestBody(body)),
	)
}

func TestUpdateCategoryName_CategoryNameNoChanged(tt *testing.T) {
	t := testing.NewTester(tt)

	body := Category{
		Name: uuid.New()[:31],
	}

	id, err := app.RegisterNewCategory(user, testing.StructToRequestBody(body))
	t.NoError(err)

	t.Is(
		service.ErrCategoryNoChanged,
		app.EditCategory(user, fmt.Sprint(id), testing.StructToRequestBody(body)),
	)
}

func TestRemoveCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)

	id, err := factory.Default().GenerateID()
	t.NoError(err)

	category, err := model.NewCategory(id, uuid.New()[:31])
	t.NoError(err)
	t.NotNil(category)

	repo := infra.NewCategoryStorage(testing.DB())
	t.NoError(repo.Add(category))

	t.NoError(app.RemoveCategoryByID(fmt.Sprint(id)))
}

func TestRemoveCategory_InvalidID(tt *testing.T) {
	t := testing.NewTester(tt)

	t.Is(service.ErrInvalidCategoryID, app.RemoveCategoryByID("invalid id !?"))
}

func TestRemoveCategory_NoSuchCategory(tt *testing.T) {
	t := testing.NewTester(tt)

	id, _ := factory.Default().GenerateID()
	err := app.RemoveCategoryByID(fmt.Sprint(id))
	t.Is(service.ErrNoSuchCategory, err)
}
