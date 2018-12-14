package service

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"lmm/api/service/asset/domain/model"
	"sync"
)

var (
	DefaultImageEncoder ImageEncoder
)

func init() {
	DefaultImageEncoder = &imageEncoder{}
}

type ImageEncoder interface {
	Encode(c context.Context, src image.Image) (model.Data, error)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type imageEncoder struct{}

func (e *imageEncoder) Encode(ctx context.Context, src image.Image) (model.Data, error) {
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	if err := jpeg.Encode(buf, src, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}

	b := buf.Bytes()[:]

	return b, nil
}
