package model

import "lmm/api/domain/model"

// Author is the model of article author
type Author struct {
	model.Entity
	name string
}

// NewAuthor returns a new author
func NewAuthor(name string) *Author {
	return &Author{name: name}
}

// Name returns the name of the author
func (a *Author) Name() string {
	return a.name
}
