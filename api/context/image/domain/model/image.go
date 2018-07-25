package model

import "time"

type Image struct {
	Entity
	id        string
	userID    uint64
	createdAt time.Time
}

func NewImage(id string, userID uint64, createdAt time.Time) *Image {
	return &Image{
		id:        id,
		userID:    userID,
		createdAt: createdAt,
	}
}

func (i *Image) ID() string {
	return i.id
}

func (i *Image) UserID() uint64 {
	return i.userID
}

func (i *Image) CreatedAt() time.Time {
	return i.createdAt
}
