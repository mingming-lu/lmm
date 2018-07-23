package ui

import (
	"lmm/api/context/blog/domain/model"
)

type Tag struct {
	Name string `json:"name"`
}

type TagListEntry struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Tags []TagListEntry

func tagsToJSON(models []*model.Tag) Tags {
	tags := make(Tags, len(models))
	for index, model := range models {
		tags[index].ID = model.ID()
		tags[index].Name = model.Name()
	}
	return tags
}
