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
	photos []*PhotoDescriptor
	nextID *uint
}

// NewPhotoCollection creates a new image collection
func NewPhotoCollection(photos []*PhotoDescriptor, nextID *uint) *PhotoCollection {
	if len(photos) == 0 {
		photos = make([]*PhotoDescriptor, 0)
	}
	return &PhotoCollection{photos: photos, nextID: nextID}
}

// List lists all photos from collection
func (c *PhotoCollection) List() []*PhotoDescriptor {
	return c.photos
}

// NextID returns next id when to find more images
// if no more images, returns nil
func (c *PhotoCollection) NextID() *uint {
	return c.nextID
}
