package model

// ImageDescriptor models the image descriptor
type ImageDescriptor struct {
	assetDescriptor
	id uint
}

// NewImageDescriptor creates new image descriptor
func NewImageDescriptor(id uint, name string) *ImageDescriptor {
	return &ImageDescriptor{
		id:              id,
		assetDescriptor: assetDescriptor{name: name},
	}
}

// ID gets id
func (d *ImageDescriptor) ID() uint {
	return d.id
}

// ImageCollection saves image
type ImageCollection struct {
	images []*ImageDescriptor
	nextID *uint
}

// NewImageCollection creates a new image collection
func NewImageCollection(images []*ImageDescriptor, nextID *uint) *ImageCollection {
	if len(images) == 0 {
		images = make([]*ImageDescriptor, 0)
	}
	return &ImageCollection{images: images, nextID: nextID}
}

// List get all images from collection
func (c *ImageCollection) List() []*ImageDescriptor {
	return c.images
}

// NextID returns next id when to find more images
// if no more images, returns nil
func (c *ImageCollection) NextID() *uint {
	return c.nextID
}
