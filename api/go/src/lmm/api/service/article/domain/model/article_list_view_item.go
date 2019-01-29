package model

import (
	"time"

	"lmm/api/domain/model"

	"github.com/pkg/errors"
)

// ArticleListViewItem is the item struct of AriticleListView
type ArticleListViewItem struct {
	model.ValueObject
	id     *ArticleID
	title  string
	postAt time.Time
}

// NewArticleListViewItem creates a new item ArticleListViewItem
func NewArticleListViewItem(rawArticleID string, title string, postAt time.Time) (*ArticleListViewItem, error) {
	articleID, err := NewArticleID(rawArticleID)
	if err != nil {
		return nil, errors.Wrap(err, rawArticleID)
	}

	return &ArticleListViewItem{id: articleID, title: title, postAt: postAt}, nil
}

// ID returns article's id
func (i *ArticleListViewItem) ID() *ArticleID {
	return i.id
}

// Title gets article's title
func (i *ArticleListViewItem) Title() string {
	return i.title
}

// PostAt gets article's post time
func (i *ArticleListViewItem) PostAt() time.Time {
	return i.postAt
}
