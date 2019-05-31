package model

// ArticleListView is the model used to view article list
type ArticleListView struct {
	items       []*ArticleListViewItem
	page        uint
	perPage     uint
	total       uint
	hasNextPage bool
}

// NewArticleListView constructs a new ArticleListView
func NewArticleListView(items []*ArticleListViewItem, page, perPage, total uint, hasNextPage bool) *ArticleListView {
	return &ArticleListView{
		items:       items,
		page:        page,
		perPage:     perPage,
		total:       total,
		hasNextPage: hasNextPage,
	}
}

// Items gets items of article list view
func (v *ArticleListView) Items() []*ArticleListViewItem {
	return v.items
}

func (v *ArticleListView) Page() uint {
	return v.page
}

func (v *ArticleListView) PerPage() uint {
	return v.perPage
}

// Total returns total articles number
func (v *ArticleListView) Total() uint {
	return v.total
}

// HasNextPage returns true if there is next page
func (v *ArticleListView) HasNextPage() bool {
	return v.hasNextPage
}
