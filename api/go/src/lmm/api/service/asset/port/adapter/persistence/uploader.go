package persistence

import (
	"context"
	"fmt"
	"io"

	"lmm/api/service/asset/usecase"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type GCSUploader struct {
	bucketName string
	bucket     *storage.BucketHandle
}

const templatePublicURL = "https://storage.googleapis.com/%s/%s"

func NewGCSUploader(c context.Context, bh *storage.BucketHandle) (*GCSUploader, error) {
	attr, err := bh.Attrs(c)
	if err != nil {
		return nil, err
	}

	return &GCSUploader{
		bucketName: attr.Name,
		bucket:     bh,
	}, nil
}

func (uploader *GCSUploader) Upload(c context.Context, asset *usecase.AssetToUpload) (string, error) {
	w := uploader.bucket.Object(asset.Filename).NewWriter(c)
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, asset.DataSource); err != nil {
		return "", errors.Wrap(err, "failed to upload to GCS")
	}

	if err := w.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close storage writer")
	}

	if err := asset.DataSource.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close file reader")
	}

	return fmt.Sprintf(templatePublicURL, uploader.bucketName, asset.Filename), nil
}
