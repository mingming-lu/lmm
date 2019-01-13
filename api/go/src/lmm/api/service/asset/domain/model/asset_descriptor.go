package model

import "lmm/api/domain/model"

type AssetDescriptor struct {
	model.ValueObject
	name      string
	assetType AssetType
}

func NewAssetDescriptor(name, assetType string) *AssetDescriptor {
	a := &AssetDescriptor{name: name}

	switch assetType {
	case "image":
		a.assetType = Image
	case "photo":
		a.assetType = Photo
	default:
		a.assetType = Unknown
	}

	return a
}

func (d *AssetDescriptor) Name() string {
	return d.name
}

func (d *AssetDescriptor) Type() AssetType {
	return d.assetType
}
