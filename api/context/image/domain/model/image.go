package model

import "time"

type Image struct {
	Entity
	id        string
	userID    uint64
	createdAt time.Time
}

func (i *Image) ID() string {
	return i.id
}

func NewImage(id string, userID uint64, createdAt time.Time) *Image {
	return &Image{
		id:        id,
		userID:    userID,
		createdAt: createdAt,
	}
}
