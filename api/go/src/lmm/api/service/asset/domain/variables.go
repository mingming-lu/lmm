package domain

import "github.com/pkg/errors"

var (
	// ErrNoSuchAsset error
	ErrNoSuchAsset = errors.New("no such asset")

	// ErrNoSuchUser error
	ErrNoSuchUser = errors.New("no such user")

	// ErrUnsupportedAssetType error
	ErrUnsupportedAssetType = errors.New("unsupported asset type")

	// ErrUnsupportedImageFormat error
	ErrUnsupportedImageFormat = errors.New("unsupported image type")

	// ErrInvalidTypeNotAPhoto error
	ErrInvalidTypeNotAPhoto = errors.New("not a photo")
)
