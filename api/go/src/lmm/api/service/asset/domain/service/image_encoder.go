package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"regexp"
	"sync"

	"github.com/pkg/errors"

	"lmm/api/service/asset/domain"
	"lmm/api/service/asset/domain/model"
)

var (
	// DefaultImageEncoder implements ImageEncoder as default
	DefaultImageEncoder ImageEncoder
)

var (
	// ErrEmptyImage occurs when read 0 byte from image file
	ErrEmptyImage = errors.New("empty image")

	// ErrMemoryOverflows error
	ErrMemoryOverflows = errors.New("memory overflows")
)

func init() {
	DefaultImageEncoder = &imageEncoder{}
}

// ImageEncoder interface
type ImageEncoder interface {
	// Encode encodes multipart.File into []byte followed by its extention
	Encode(c context.Context, file multipart.File) (bytes model.Data, extention string, err error)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type imageEncoder struct{}

func (e *imageEncoder) Encode(ctx context.Context, file multipart.File) (data model.Data, ext string, err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		panicErr, ok := e.(error)
		if !ok {
			panic(e)
		}
		switch panicErr {
		case bytes.ErrTooLarge: // panicked by bytes.Buffer.ReadFrom
			err = ErrMemoryOverflows
		default:
			panic(e)
		}
	}()

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	var n int64
	n, err = buf.ReadFrom(file)
	if err != nil {
		return
	}
	if n == 0 {
		err = ErrEmptyImage
		return
	}

	data = buf.Bytes()
	ext, err = assignImageExtension(data)

	return
}

var contentTypePattern = regexp.MustCompile(`^(.+)\/(.+)$`)

func assignImageExtension(data []byte) (string, error) {
	ct := http.DetectContentType(data)
	switch ct {
	case "image/bmp":
		return "bmp", nil
	case "image/gif":
		return "gif", nil
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	case "image/webp":
		return "webp", nil
	default:
		ss := contentTypePattern.FindStringSubmatch(ct)
		if len(ss) != 3 {
			return ct, domain.ErrUnsupportedImageFormat
		}
		return ss[2], domain.ErrUnsupportedImageFormat
	}
}
