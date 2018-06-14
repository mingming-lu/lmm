package model

import (
	"errors"
	"lmm/api/domain/model"
)

var (
	ErrInvalidCategoryName = errors.New("invalid category name")
)

type Category struct {
	model.Entity
	id   uint64
	name string
}

func NewCategory(id uint64, name string) (*Category, error) {
	c := &Category{
		id:   id,
		name: name,
	}
	if c.validateName(c.name) != nil {
		return nil, ErrInvalidCategoryName
	}
	return c, nil
}

func (c *Category) ID() uint64 {
	return c.id
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) UpdateName(newName string) error {
	if c.validateName(newName) != nil {
		return ErrInvalidCategoryName
	}
	c.name = newName
	return nil
}

func (c *Category) validateName(name string) error {
	// TODO validating name by regexp
	return nil
}
