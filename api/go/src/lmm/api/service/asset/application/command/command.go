package command

import (
	"mime/multipart"

	"lmm/api/service/asset/domain/model"
)

// UploadAsset command
type UploadAsset struct {
	userID    string
	assetType string
	file      multipart.File
}

func NewUploadAsset(userID string, assetType string, file multipart.File) *UploadAsset {
	return &UploadAsset{
		userID:    userID,
		assetType: assetType,
		file:      file,
	}
}

// UserID getter
func (cmd *UploadAsset) UserID() string {
	return cmd.userID
}

// Type adapts asset type string to asset type model
func (cmd *UploadAsset) Type() model.AssetType {
	switch cmd.assetType {
	case "image":
		return model.Image
	case "photo":
		return model.Photo
	default:
		return &unknownAssetType{s: cmd.assetType}
	}
}

// File gets mutipart.File
func (cmd *UploadAsset) File() multipart.File {
	return cmd.file
}

type unknownAssetType struct {
	s string
}

func (t unknownAssetType) String() string {
	return t.s
}
