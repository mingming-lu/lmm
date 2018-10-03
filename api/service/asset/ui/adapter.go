package ui

import "lmm/api/service/asset/domain/model"

type imageListItem struct {
	Name string `json:"name"`
}

type imageListJSON struct {
	Images      []imageListItem `json:"images"`
	HasNextPage bool            `json:"hasNextPage"`
}

func imageCollectionToJSON(collection *model.ImageCollection) *imageListJSON {
	images := make([]imageListItem, len(collection.List()), len(collection.List()))
	for i, image := range collection.List() {
		images[i].Name = image.Name()
	}
	return &imageListJSON{
		Images:      images,
		HasNextPage: collection.HasNextPage(),
	}
}

type photoListItem struct {
	Name string `json:"name"`
}

type photoListJSON struct {
	Photos      []photoListItem `json:"photos"`
	HasNextPage bool            `json:"hasNextPage"`
}

func photoCollectionToJSON(collection *model.PhotoCollection) *photoListJSON {
	photos := make([]photoListItem, len(collection.List()), len(collection.List()))
	for i, photo := range collection.List() {
		photos[i].Name = photo.Name()
	}
	return &photoListJSON{
		Photos:      photos,
		HasNextPage: collection.HasNextPage(),
	}
}
