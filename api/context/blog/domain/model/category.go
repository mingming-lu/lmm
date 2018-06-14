package model

import "lmm/api/domain/model"

type Category struct {
	model.Entity
	id   uint64
	name string
}

func (c *Category) ID() uint64 {
	return c.id
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) UpdateCategory(newName string) {
	c.name = newName
}
