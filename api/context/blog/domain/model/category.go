package model

import (
	"errors"
	"lmm/api/domain/model"
	"lmm/api/utils/strings"
	"regexp"
)

var (
	ErrInvalidCategoryName = errors.New("invalid category name")
)

var (
	patternValidCategoryName = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]{1,31}$")
)

type Category struct {
	model.Entity
	id   uint64
	name string
}

func NewCategory(id uint64, name string) (*Category, error) {
	name, err := validateCategoryName(name)
	if err != nil {
		return nil, ErrInvalidCategoryName
	}

	c := &Category{
		id:   id,
		name: name,
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
	newName, err := validateCategoryName(newName)
	if err != nil {
		return ErrInvalidCategoryName
	}
	c.name = newName
	return nil
}

func validateCategoryName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if !patternValidCategoryName.MatchString(name) {
		return "", ErrInvalidCategoryName
	}
	return name, nil
}
