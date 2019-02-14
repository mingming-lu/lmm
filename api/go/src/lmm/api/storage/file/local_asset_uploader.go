package file

import (
	"bufio"
	"context"
	"os"

	"github.com/pkg/errors"
)

var (
	// ErrImageUploadTypeNotGiven error
	ErrImageUploadTypeNotGiven = errors.New("asset upload type not given")

	// ErrInvalidImageUploadType error
	ErrInvalidImageUploadType = errors.New("asset type expects to be 'image' or 'photo'")

	ErrInvalidAssetName = errors.New("invalid asset name")
)

// ImageUploaderConfig pass config parameters to uploader
type ImageUploaderConfig struct {
	Type string
}

// LocalImageUploader implements Uploader for uploading images to server's local disk
// Not support sharding
type localImageUploader struct {
	basePath string
}

// NewLocalImageUploader returns a Uploader implementation
func NewLocalImageUploader() Uploader {
	return &localImageUploader{basePath: "/asset/"}
}

func (up *localImageUploader) Upload(c context.Context, name string, data []byte, args ...interface{}) error {
	if len(args) != 1 {
		return ErrImageUploadTypeNotGiven
	}
	if config, ok := args[0].(ImageUploaderConfig); ok {
		switch config.Type {
		case "image", "photo":
			name = up.basePath + config.Type + "s/" + name
		default:
			return errors.Wrap(ErrInvalidImageUploadType, "invalid asset upload type: "+config.Type)
		}
	} else {
		return ErrImageUploadTypeNotGiven
	}

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
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

func (up *localImageUploader) Delete(c context.Context, name string, args ...interface{}) error {
	return errors.New("not implemented")
}

func (up *localImageUploader) Close() error {
	return nil
}
