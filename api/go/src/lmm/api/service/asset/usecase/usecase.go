package usecase

import (
	"context"
	"io"
	"path"
	"time"

	"lmm/api/clock"
	"lmm/api/pkg/transaction"
	"lmm/api/util/stringutil"
	"lmm/api/util/uuidutil"

	"github.com/pkg/errors"
)

type Asset struct {
	Filename   string
	Type       string
	UploadedAt time.Time
	UserID     int64
}

type AssetToUpload struct {
	ContentType string
	DataSource  io.ReadCloser
	Filename    string
	UserID      int64
}

type FileUploader interface {
	Upload(c context.Context, assert *AssetToUpload) (string, error)
}

type Photo struct {
	URL string `json:"url"`
}

type AssetRepository interface {
	Save(c context.Context, asset *Asset) error
	ListPhotos(c context.Context, count int, cursor string) ([]*Photo, string, error)
}

type Usecase struct {
	assetRepository AssetRepository
	fileUploader    FileUploader
	txManager       transaction.Manager
}

func New(assertRepository AssetRepository, fileUploader FileUploader, txManager transaction.Manager) *Usecase {
	return &Usecase{
		assetRepository: assertRepository,
		fileUploader:    fileUploader,
		txManager:       txManager,
	}
}

func (uc *Usecase) UploadPhoto(c context.Context, photo *AssetToUpload) (url string, err error) {
	// rename photo filename randomly
	if photo.Filename == "" {
		panic("internal error: empty filename")
	}
	photo.Filename = uuidutil.NewUUID() + path.Ext(photo.Filename)

	err = uc.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		if err := uc.assetRepository.Save(tx, &Asset{
			Filename:   photo.Filename,
			UploadedAt: clock.Now(),
			UserID:     photo.UserID,
			Type:       "Photo",
		}); err != nil {
			return err
		}

		url, err = uc.fileUploader.Upload(tx, photo)
		if err != nil {
			return errors.Wrap(err, "failed to upload photo")
		}

		return err
	}, nil)

	return
}

func (uc *Usecase) UploadAsset(c context.Context, assert *AssetToUpload) error {
	panic("not implemented")
}

func (uc *Usecase) ListPhotos(c context.Context, countStr, cursor string) ([]*Photo, string, error) {
	count, err := stringutil.ParseInt(countStr)
	if err != nil || count < 1 {
		return nil, "", errors.New("invalid count")
	}

	return uc.assetRepository.ListPhotos(c, count, cursor)
}
