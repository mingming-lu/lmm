package model

import "lmm/api/domain/model"

// ArticleListView is the model used to view article list
type ArticleListView struct {
	model.ValueObject
	items       []*ArticleListViewItem
	hasNextPage bool
}

// NewArticleListView constructs a new ArticleListView
func NewArticleListView(items []*ArticleListViewItem, hasNextPage bool) *ArticleListView {
	return &ArticleListView{items: items, hasNextPage: hasNextPage}
}

// Items gets items of article list view
func (v *ArticleListView) Items() []*ArticleListViewItem {
	return v.items
}

// HasNextPage returns true if there is next page
func (v *ArticleListView) HasNextPage() bool {
	return v.hasNextPage
}
