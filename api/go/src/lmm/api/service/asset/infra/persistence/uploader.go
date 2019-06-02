package persistence

import (
	"context"
	"io"
	"os"

	"lmm/api/service/asset/usecase"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

type GCSUploader struct {
	gcsClient *storage.Client
	bucket    *storage.BucketHandle
}

var (
	bucketName    = os.Getenv("PROJECT_ID") + "-asset"
	publicURLBase = "https://storage.googleapis.com/" + bucketName + "/"
)

func NewGCSUploader(client *storage.Client) *GCSUploader {
	bucket := client.Bucket(bucketName)

	return &GCSUploader{
		gcsClient: client,
		bucket:    bucket,
	}
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

	return publicURLBase + asset.Filename, nil
}
