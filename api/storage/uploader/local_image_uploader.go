package uploader

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

// LocalImageUploader implements Uploader for uploading images to server's local disk
// Not support sharding
type localImageUploader struct {
	path string
}

// NewLocalImageUploader returns a Uploader implementation
func NewLocalImageUploader() Uploader {
	path := "/static/images/"
	if err := os.MkdirAll(path, 0644); err != nil {
		panic(err)
	}
	return &localImageUploader{path: path}
}

func (up *localImageUploader) Upload(name string, data []byte, args ...interface{}) error {
	file, err := os.OpenFile(up.path+name, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	if _, err := w.Write(data); err != nil {
		return err
	}

	return w.Flush()
}

func (up *localImageUploader) Delete(name string, args ...interface{}) error {
	return errors.New("not implemented")
}
