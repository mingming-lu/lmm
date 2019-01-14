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

// SetImageAlternateTexts is the command to set image's alt
type SetImageAlternateTexts struct {
	imageName string
	altNames  []string
}

// NewSetImageAlternateTexts creates a new SetImageAlternateTexts command
func NewSetImageAlternateTexts(imageName string, altNames []string) *SetImageAlternateTexts {
	return &SetImageAlternateTexts{
		imageName: imageName,
		altNames:  altNames,
	}
}

// ImageName returns image's name
func (cmd *SetImageAlternateTexts) ImageName() string {
	return cmd.imageName
}

// AltNames returns alt names to set
func (cmd *SetImageAlternateTexts) AltNames() []string {
	return cmd.altNames
}
