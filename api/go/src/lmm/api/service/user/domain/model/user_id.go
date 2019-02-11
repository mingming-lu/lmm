package model

import (
	"regexp"

	"lmm/api/service/base/model"
	"lmm/api/service/user/domain"
)

var (
	patternUserID = regexp.MustCompile(`^\d{8}$`)
)

// UserID domain model
type UserID struct {
	model.ValueObject
	id string
}

// NewUserID builds a new UserID pointer from given raw id
func NewUserID(rawID string) *UserID {
	return &UserID{id: rawID}
}

func (id *UserID) setID(anID string) error {
	if !patternUserID.MatchString(anID) {
		return domain.ErrInvalidUserID
	}
	id.id = anID
	return nil
}

func (id *UserID) String() string {
	return id.id
}
