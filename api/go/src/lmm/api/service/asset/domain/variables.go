package domain

import "github.com/pkg/errors"

var (
	// ErrNoSuchUser error
	ErrNoSuchUser = errors.New("no such user")

	// ErrUnsupportedAssetType error
	ErrUnsupportedAssetType = errors.New("unsupported asset type")

	// ErrUnsupportedImageFormat error
	ErrUnsupportedImageFormat = errors.New("unsupported image type")
)
