package model

// ImageDescriptor models the image descriptor
type ImageDescriptor struct {
	AssetDescriptor
	id uint
}

// NewImageDescriptor creates new image descriptor
func NewImageDescriptor(id uint, name string) *ImageDescriptor {
	return &ImageDescriptor{
		id:              id,
		AssetDescriptor: AssetDescriptor{name: name, assetType: Image},
	}
}

// ID gets id
func (d *ImageDescriptor) ID() uint {
	return d.id
}

// ImageCollection saves image
type ImageCollection struct {
	images      []*ImageDescriptor
	hasNextPage bool
}

// NewImageCollection creates a new image collection
func NewImageCollection(images []*ImageDescriptor, hasNextPage bool) *ImageCollection {
	if len(images) == 0 {
		images = make([]*ImageDescriptor, 0)
	}
	return &ImageCollection{images: images, hasNextPage: hasNextPage}
}

// List get all images from collection
func (c *ImageCollection) List() []*ImageDescriptor {
	return c.images
}

// HasNextPage returns next id when to find more images
// if no more images, returns nil
func (c *ImageCollection) HasNextPage() bool {
	return c.hasNextPage
}
