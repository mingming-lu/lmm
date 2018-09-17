package model

import "lmm/api/context/base/domain/model"

// Author is the model of article author
type Author struct {
	model.Entity
	id   int64
	name string
}

// NewAuthor returns a new author
func NewAuthor(id int64, name string) Author {
	return Author{id: id, name: name}
}

// ID returns the id of the author
func (a Author) ID() int64 {
	return a.id
}

// Name returns the name of the author
func (a Author) Name() string {
	return a.name
}
