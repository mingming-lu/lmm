package model

import (
	"time"
)

// ArticleListViewItem is the item struct of AriticleListView
type ArticleListViewItem struct {
	id       int64
	linkName string
	title    string
	postAt   time.Time
}

// NewArticleListViewItem creates a new item ArticleListViewItem
func NewArticleListViewItem(id int64, linkName, title string, postAt time.Time) (*ArticleListViewItem, error) {
	return &ArticleListViewItem{
		id:       id,
		linkName: linkName,
		title:    title,
		postAt:   postAt,
	}, nil
}

// ID returns article's id
func (i *ArticleListViewItem) ID() int64 {
	return i.id
}

func (i *ArticleListViewItem) LinkName() string {
	return i.linkName
}

// Title gets article's title
func (i *ArticleListViewItem) Title() string {
	return i.title
}

// PostAt gets article's post time
func (i *ArticleListViewItem) PostAt() time.Time {
	return i.postAt
}
