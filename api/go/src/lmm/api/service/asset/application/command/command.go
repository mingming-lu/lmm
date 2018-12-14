package command

import "mime/multipart"

type UploadAsset struct {
	username  string
	extention string
	file      multipart.File
}

func (c *UploadAsset) UploaderName() string {
	return c.username
}

func (c *UploadAsset) Extention() string {
	return c.extention
}

func (c *UploadAsset) File() multipart.File {
	return c.file
}
