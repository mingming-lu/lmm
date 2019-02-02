package model

import "lmm/api/model"

type AssetDescriptor struct {
	model.ValueObject
	name      string
	assetType AssetType
}

func NewAssetDescriptor(name string, assetType AssetType) *AssetDescriptor {
	a := &AssetDescriptor{name: name}

	switch assetType {
	case Image:
		a.assetType = Image
	case Photo:
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
