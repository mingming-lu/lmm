package model

import "lmm/api/domain/model"

type Tag struct {
	model.Entity
	name string
}
