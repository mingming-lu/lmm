package domain

import "github.com/pkg/errors"

var (
	// ErrNoSuchAsset error
	ErrNoSuchAsset = errors.New("no such asset")

	// ErrNoSuchUser error
	ErrNoSuchUser = errors.New("no such user")

	// ErrNoSuchPhoto error
	ErrNoSuchPhoto = errors.New("no such photo")

	// ErrUnsupportedAssetType error
	ErrUnsupportedAssetType = errors.New("unsupported asset type")

	// ErrUnsupportedImageFormat error
	ErrUnsupportedImageFormat = errors.New("unsupported image type")

	// ErrInvalidTypeNotAPhoto error
	ErrInvalidTypeNotAPhoto = errors.New("not a photo")

	// ErrDuplicateImageAlt error
	ErrDuplicateImageAlt = errors.New("duplicate image alternate text")
)
