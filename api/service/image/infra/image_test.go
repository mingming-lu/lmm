package infra

import (
	"lmm/api/context/image/domain"
	"lmm/api/context/image/domain/factory"
	"lmm/api/storage/static"
	"lmm/api/testing"
)

func TestAddImage_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)

	t.NoError(repo.Add(model.WrapData(nil)))
}

func TestAddImage_Duplicate(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)

	t.NoError(repo.Add(model.WrapData(nil)))
	t.IsError(domain.ErrDuplicateImageID, repo.Add(model.WrapData(nil)))
}

func RemoveImage_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)
	t.NoError(repo.Add(model.WrapData(nil)))
	t.NoError(repo.Remove(model))
}

func TestRemoveImage_NoSuchImage(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)
	t.IsError(domain.ErrNoSuchImage, repo.Remove(model))
}

func TestFindImageByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)
	t.NoError(repo.Add(model.WrapData(nil)))

	modelFound, err := repo.FindByID(model.ID())
	t.NoError(err)
	t.Is(model.ID(), modelFound.ID())
	t.Is(model.UserID(), modelFound.UserID())
}

func TestFindImageByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)
	modelFound, err := repo.FindByID(model.ID())
	t.IsError(domain.ErrNoSuchImage, err)
	t.Nil(modelFound)
}

func TestAddImage_StaticFile(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())
	repo.SetStaticRepository(&static.LocalStaticRepository{})

	model := factory.NewImage(1)

	t.NoError(repo.Add(model.WrapData(nil)))
	t.NoError(repo.Remove(model))
}
