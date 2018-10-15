package model

import "fmt"

// AssetType is a type
type AssetType fmt.Stringer

// Image asset type
var Image AssetType = &image{}

type image struct{}

func (i image) String() string {
	return "image"
}

// Photo asset type
var Photo AssetType = &photo{}

type photo struct{}

func (p photo) String() string {
	return "photo"
}
