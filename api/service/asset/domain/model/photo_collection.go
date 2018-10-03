package model

// PhotoDescriptor models the image descriptor
type PhotoDescriptor struct {
	assetDescriptor
	id uint
}

// NewPhotoDescriptor creates a new photo descriptor
func NewPhotoDescriptor(id uint, name string) *PhotoDescriptor {
	return &PhotoDescriptor{
		id:              id,
		assetDescriptor: assetDescriptor{name: name},
	}
}

// ID gets id
func (d *PhotoDescriptor) ID() uint {
	return d.id
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
