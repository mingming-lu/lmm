package model

// PhotoDescriptor models the image descriptor
type PhotoDescriptor struct {
	assetDescriptor
}

// NewPhotoDescriptor creates a new photo descriptor
func NewPhotoDescriptor(name string) *PhotoDescriptor {
	return &PhotoDescriptor{assetDescriptor: assetDescriptor{name: name}}
}

// PhotoCollection saves photos
type PhotoCollection struct {
	photos []*PhotoDescriptor
}

// List lists all photos from collection
func (c *PhotoCollection) List() []*PhotoDescriptor {
	return c.photos
}
