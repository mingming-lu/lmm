package model

import "lmm/api/domain/model"

type Tag struct {
	model.Entity
	id   uint64
	data TagData
}

type TagData struct {
	model.ValueObject
	blogID uint64
	name   string
}

func NewTag(id, blogID uint64, name string) (*Tag, error) {
	data := TagData{blogID: blogID, name: name}
	return &Tag{id: id, data: data}, nil
}

func (t *Tag) ID() uint64 {
	return t.id
}

func (t *Tag) BlogID() uint64 {
	return t.data.blogID
}

func (t *Tag) Name() string {
	return t.data.name
}

func (t *Tag) UpdateName(newName string) error {
	// TODO validate name
	t.data.name = newName
	return nil
}
