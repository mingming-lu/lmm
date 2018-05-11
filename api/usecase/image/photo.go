package image

import (
	model "lmm/api/domain/model/image"
	repo "lmm/api/domain/repository/image"
	"strconv"

	"github.com/akinaru-lu/errors"
)

var (
	ErrInvalidUserID = errors.New("Invalid user ID")
	ErrInvalidCount  = errors.New("Invalid count")
	ErrInvalidPage   = errors.New("Invalid page")
)

type FetchPhotosResponse struct {
	Images  []model.Minimal `json:"images"`
	HasNext bool            `json:"has_next"`
}

func FetchPhotos(userIDStr, countStr, pageStr string) (*FetchPhotosResponse, error) {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	count := uint64(10)
	if countStr != "" {
		_count, err := strconv.ParseUint(countStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidCount
		}
		count = _count
	}

	page := uint64(1)
	if pageStr != "" {
		_page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			return nil, ErrInvalidPage
		}
		page = _page
	}

	urls, err := repo.SearchPhotosByUserID(userID, count+1, page)
	hasNext := false
	if uint64(len(urls)) > count {
		hasNext = true
	}
	res := &FetchPhotosResponse{
		Images:  urls,
		HasNext: hasNext,
	}
	return res, err
}

func TurnOnPhotoSwitch(userID int64, imageName string) error {
	return togglePhotoSwitch(userID, imageName, true)
}

func TurnOffPhotoSwitch(userID int64, imageName string) error {
	return togglePhotoSwitch(userID, imageName, false)
}

func togglePhotoSwitch(userID int64, imageName string, shown bool) error {
	image, err := repo.ByName(userID, imageName)
	if err != nil {
		return err
	}
	return repo.SavePhoto(userID, image.ID, shown)
}
