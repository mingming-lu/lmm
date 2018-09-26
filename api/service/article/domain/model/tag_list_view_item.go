package model

import "lmm/api/domain/model"

// TagListViewItem is a model used to view tag
type TagListViewItem struct {
	model.ValueObject
	name string
}

// NewTagListViewItem creates a new TagListViewItem
func NewTagListViewItem(name string) *TagListViewItem {
	return &TagListViewItem{name: name}
}

// Name gets its name
func (i *TagListViewItem) Name() string {
	return i.name
}
