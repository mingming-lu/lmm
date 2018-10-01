package model

import (
	"lmm/api/domain/model"
)

// Image domain entity
type Image struct {
	model.Entity
	name     string
	uploader *Uploader
	data     Data
}

// NewImage creates new image entity
func NewImage(name string, uploader *Uploader, data []byte) *Image {
	return &Image{name: name, uploader: uploader, data: data}
}

// Name gets image's name
func (i *Image) Name() string {
	return i.name
}

// Uploader gets image's uploader
func (i *Image) Uploader() *Uploader {
	return i.uploader
}

// Data gets image's data
func (i *Image) Data() []byte {
	return i.data
}
