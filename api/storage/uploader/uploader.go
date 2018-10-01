package uploader

// Uploader provies interface to deal with uploading
type Uploader interface {
	Upload(id string, data []byte, args ...interface{}) error
	Delete(id string, args ...interface{}) error
}
