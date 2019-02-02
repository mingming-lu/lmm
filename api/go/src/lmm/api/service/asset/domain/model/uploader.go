package model

import "lmm/api/model"

// Uploader models image's uploader
type Uploader struct {
	model.Entity
	id UploaderID
}

// NewUploader creates a new Uploader pointer
func NewUploader(id UploaderID) *Uploader {
	return &Uploader{id: id}
}

// ID returns uploader's id
func (uploader *Uploader) ID() UploaderID {
	return uploader.id
}
