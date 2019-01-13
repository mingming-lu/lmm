package model

// Alt is the alt of image
type Alt struct {
	id *altID
}

type altID struct {
	assetName string
	name      string
}

// NewAlt creates a new alt model
func NewAlt(assetName, altName string) *Alt {
	return &Alt{
		id: &altID{assetName: assetName, name: altName},
	}
}

// Name gets alt's name
func (alt *Alt) Name() string {
	return alt.id.name
}
