package model

import (
	"lmm/api/service/asset/domain"
	"lmm/api/util/sliceutil"
)

// PhotoDescriptor models the image descriptor
type PhotoDescriptor struct {
	AssetDescriptor
	id   uint
	alts []string
}

// NewPhotoDescriptor creates a new photo descriptor
func NewPhotoDescriptor(id uint, name string) *PhotoDescriptor {
	return &PhotoDescriptor{
		id:              id,
		alts:            make([]string, 0),
		AssetDescriptor: AssetDescriptor{name: name, assetType: Photo},
	}
}

// ID gets id
func (d *PhotoDescriptor) ID() uint {
	return d.id
}

// AddAlternateText appends text into d's alternate texts list
func (d *PhotoDescriptor) AddAlternateText(text string) error {
	if sliceutil.ContainsString(text, d.alts) {
		return domain.ErrDuplicateImageAlt
	}
	d.alts = append(d.alts, text)
	return nil
}

// AlternateTexts returns photo's alts
func (d *PhotoDescriptor) AlternateTexts() []string {
	return d.alts
}

// PhotoCollection saves photos
type PhotoCollection struct {
	photos      []*PhotoDescriptor
	hasNextPage bool
}

// NewPhotoCollection creates a new image collection
func NewPhotoCollection(photos []*PhotoDescriptor, hasNextPage bool) *PhotoCollection {
	if len(photos) == 0 {
		photos = make([]*PhotoDescriptor, 0)
	}
	return &PhotoCollection{photos: photos, hasNextPage: hasNextPage}
}

// List lists all photos from collection
func (c *PhotoCollection) List() []*PhotoDescriptor {
	return c.photos
}

// HasNextPage returns next id when to find more images
// if no more images, returns nil
func (c *PhotoCollection) HasNextPage() bool {
	return c.hasNextPage
}
