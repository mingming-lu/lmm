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

type AssetType string

func (t AssetType) String() string {
	return string(t)
}

func AssetTypeFromString(name string) AssetType {
	switch name {
	case "Image":
		return ImageType
	case "Photo":
		return PhotoType
	default:
		return UnknownType
	}
}

const (
	ImageType   AssetType = "Image"
	PhotoType   AssetType = "Photo"
	UnknownType AssetType = "Unknown"
)

var (
	ErrNoSuchPhoto = errors.New("no such photo")
	ErrNotPhoto    = errors.New("not a photo")
)

type Asset struct {
	ID         *AssetID
	Filename   string
	Type       AssetType
	UploadedAt time.Time
}

type AssetID struct {
	ID     int64
	UserID int64
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
	URL  string   `json:"url"`
	Tags []string `json:"tags"`
}

type AssetRepository interface {
	NextID(c context.Context, userID int64) (*AssetID, error)
	Save(c context.Context, asset *Asset) error
	FindByFileName(c context.Context, filename string) (*Asset, error)
	SetPhotoTags(c context.Context, id *AssetID, tags []string) error
	ListPhotos(c context.Context, count int, cursor string) ([]*Photo, string, error)
	GetPublicURL(c context.Context, filename string) string
	GetTagsByPhotoID(c context.Context, id *AssetID) ([]string, error)
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

	id, err := uc.assetRepository.NextID(c, photo.UserID)
	if err != nil {
		return "", err
	}

	err = uc.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		if err := uc.assetRepository.Save(tx, &Asset{
			ID:         id,
			Filename:   photo.Filename,
			UploadedAt: clock.Now(),
			Type:       PhotoType,
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

func (uc *Usecase) SetPhotoTags(c context.Context, filename string, tags []string) error {
	return uc.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		asset, err := uc.assetRepository.FindByFileName(tx, filename)
		if err != nil {
			return errors.Wrap(ErrNoSuchPhoto, err.Error())
		}

		if asset.Type != PhotoType {
			return ErrNotPhoto
		}

		return uc.assetRepository.SetPhotoTags(tx, asset.ID, tags)
	}, nil)
}

func (uc *Usecase) UploadAsset(c context.Context, assert *AssetToUpload) error {
	panic("not implemented")
}

func (uc *Usecase) GetPhotoInfo(c context.Context, filename string) (photo *Photo, err error) {
	err = uc.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		asset, err := uc.assetRepository.FindByFileName(tx, filename)
		if asset.Type != PhotoType {
			return ErrNotPhoto
		}

		tags, err := uc.assetRepository.GetTagsByPhotoID(tx, asset.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get photo tags")
		}

		photo = &Photo{
			URL:  uc.assetRepository.GetPublicURL(tx, asset.Filename),
			Tags: tags,
		}

		return err
	}, &transaction.Option{ReadOnly: true})

	return
}

func (uc *Usecase) ListPhotos(c context.Context, countStr, cursor string) (photos []*Photo, next string, err error) {
	var count int
	count, err = stringutil.ParseInt(countStr)
	if err != nil || count < 1 {
		err = errors.Wrap(err, "invalid count")
		return
	}

	err = uc.txManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		photos, next, err = uc.assetRepository.ListPhotos(tx, count, cursor)

		return err
	}, &transaction.Option{ReadOnly: true})

	return
}
