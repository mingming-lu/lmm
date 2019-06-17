package model

type Author struct {
	id int64
}

func NewAuthor(id int64) *Author {
	return &Author{id}
}

func (a *Author) ID() int64 {
	return a.id
}
