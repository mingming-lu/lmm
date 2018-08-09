package ui

import (
	"lmm/api/context/blog/domain/model"
)

type Tag struct {
	Name string `json:"name"`
}

type TagListEntry struct {
	Name string `json:"name"`
}

type TagList struct {
	Tags []TagListEntry `json:"tags"`
}

func tagsToJSON(models []*model.Tag) TagList {
	tags := make([]TagListEntry, len(models))
	for index, model := range models {
		tags[index].Name = model.Name()
	}
	return TagList{Tags: tags}
}
