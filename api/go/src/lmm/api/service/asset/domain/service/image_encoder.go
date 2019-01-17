package service

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"lmm/api/service/asset/domain/model"
)

var (
	// ErrEmptyImage occurs when read 0 byte from image file
	ErrEmptyImage = errors.New("empty image")

	// ErrMemoryOverflows error
	ErrMemoryOverflows = errors.New("memory overflows")

	// ErrUnsupportedImageFormat error
	ErrUnsupportedImageFormat = errors.New("unsupported image format")
)

// ImageEncoder interface
type ImageEncoder interface {
	// Encode encodes multipart.File into []byte followed by its extention
	Encode(context.Context, []byte) (bytes model.Data, extention string, err error)
}

// NopImageEncoder detects extention of data and do nothing about encode
type NopImageEncoder struct{}

// Encode implementation
func (encoder *NopImageEncoder) Encode(c context.Context, data []byte) (model.Data, string, error) {
	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/gif":
		return model.Data(data), "gif", nil
	case "image/jpeg":
		return model.Data(data), "jpeg", nil
	case "image/png":
		return model.Data(data), "png", nil
	case "image/webp":
		return model.Data(data), "webp", nil
	case "image/bmp":
		return model.Data(data), "bmp", nil
	default:
		println("wtf?", string(data[:]))
		return nil, "", errors.Wrap(ErrUnsupportedImageFormat, contentType)
	}
}
