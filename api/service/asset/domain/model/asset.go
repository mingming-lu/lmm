package model

import (
	"lmm/api/domain/model"
)

// Asset domain entity
type Asset struct {
	model.Entity
	name      string
	assetType AssetType
	uploader  *Uploader
	data      Data
}

// NewAsset creates new asset entity
func NewAsset(t AssetType, name string, uploader *Uploader, data []byte) *Asset {
	return &Asset{assetType: t, name: name, uploader: uploader, data: data}
}

// Name gets asset's name
func (i *Asset) Name() string {
	return i.name
}

// Uploader gets asset's uploader
func (i *Asset) Uploader() *Uploader {
	return i.uploader
}

// Data gets asset's data
func (i *Asset) Data() []byte {
	return i.data
}

// Type gets asset's type
func (i *Asset) Type() AssetType {
	return i.assetType
}
