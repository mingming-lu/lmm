package infra

import (
	"lmm/api/context/image/domain"
	"lmm/api/context/image/domain/factory"
	"lmm/api/testing"
)

func TestAddImage_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)

	t.NoError(repo.Add(model))
}

func TestAddImage_Duplicate(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)

	t.NoError(repo.Add(model))
	t.IsError(domain.ErrDuplicateImageID, repo.Add(model))
}

func RemoveImage_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewImageStorage(testing.DB())

	model := factory.NewImage(1)
	t.NoError(repo.Add(model))
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
	t.NoError(repo.Add(model))

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
