package uploader

import "context"

// Uploader provies interface to deal with uploading
type Uploader interface {
	Upload(c context.Context, id string, data []byte, args ...interface{}) error
	Delete(c context.Context, id string, args ...interface{}) error
}
