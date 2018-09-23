package fetcher

import (
	"lmm/api/context/article/domain/model"
	"lmm/api/storage"
	"time"
)

// ArticleFetcher implements domain.model.finder.ArticleFinder
type ArticleFetcher struct {
	db *storage.DB
}

// ListByPage implementation
func (f *ArticleFetcher) ListByPage(count, page uint) (*model.ArticleListView, error) {
	stmt := f.db.MustPrepare(`select uid, title, created_at from article limit ? offset ?`)
	defer stmt.Close()

	rows, err := stmt.Query(count+1, (page-1)*count)
	if err != nil {
		if err == storage.ErrNoRows {
			return model.NewArticleListView(nil, false), nil
		}
		return nil, err
	}
	defer rows.Close()

	var (
		rawArticleID  string
		articleTitle  string
		articlePostAt time.Time
	)
	articles := make([]*model.ArticleListViewItem, 0)
	for rows.Next() {
		if err := rows.Scan(&rawArticleID, &articleTitle, &articlePostAt); err != nil {
			return nil, err
		}
		article, err := model.NewArticleListViewItem(rawArticleID, articleTitle, articlePostAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	hasNextPage := false
	if uint(len(articles)) > count {
		hasNextPage = true
		articles = articles[0:count]
	}

	return model.NewArticleListView(articles, hasNextPage), nil
}

// FindByID implementation
func (f *ArticleFetcher) FindByID(id *model.ArticleID) (*model.ArticleView, error) {
	panic("not implemented")
}
