package model

// Author is the model of article author
type Author struct {
	id int64
}

// NewAuthor returns a new author
func NewAuthor(id int64) *Author {
	return &Author{id: id}
}

// ID returns author id
func (a *Author) ID() int64 {
	return a.id
}
