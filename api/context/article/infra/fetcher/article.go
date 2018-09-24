package fetcher

import (
	"time"

	"lmm/api/context/article/domain"
	"lmm/api/context/article/domain/model"
	"lmm/api/storage"
)

// ArticleFetcher implements domain.model.finder.ArticleFinder
type ArticleFetcher struct {
	db *storage.DB
}

// NewArticleFetcher creates new ArticleFetcher
func NewArticleFetcher(db *storage.DB) *ArticleFetcher {
	return &ArticleFetcher{db: db}
}

// ListByPage implementation
func (f *ArticleFetcher) ListByPage(count, page uint) (*model.ArticleListView, error) {
	stmt := f.db.MustPrepare(`select uid, title, created_at from article order by created_at desc limit ? offset ?`)
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
	selectArticle := f.db.MustPrepare(`select id, uid, title, body, created_at, updated_at from article where uid = ?`)
	defer selectArticle.Close()

	selectTags := f.db.MustPrepare(`select sort, name from article_tag where article = ?`)
	defer selectTags.Close()

	var (
		linkedID        uint
		rawArticleID    string
		articleTitle    string
		articleBody     string
		articlePostAt   time.Time
		articleEditedAt time.Time
	)

	err := selectArticle.QueryRow(id.String()).Scan(&linkedID, &rawArticleID, &articleTitle, &articleBody, &articlePostAt, &articleEditedAt)
	if err != nil {
		if err == storage.ErrNoRows {
			return nil, domain.ErrNoSuchArticle
		}
		return nil, err
	}

	articleID, err := model.NewArticleID(rawArticleID)
	if err != nil {
		return nil, err
	}
	articleText, err := model.NewText(articleTitle, articleBody)
	if err != nil {
		return nil, err
	}

	var (
		tagOrder uint
		tagName  string
	)

	rows, err := selectTags.Query(linkedID)
	if err != nil && err != storage.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*model.Tag, 0)
	for rows.Next() {
		if err := rows.Scan(&tagOrder, &tagName); err != nil {
			return nil, err
		}
		tag, err := model.NewTag(articleID, tagOrder, tagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	content, err := model.NewContent(articleText, tags)
	if err != nil {
		return nil, err
	}

	return model.NewArticleView(articleID, content, articlePostAt, articleEditedAt), nil
}

// ListAllTags implementation
func (f *ArticleFetcher) ListAllTags() (model.TagListView, error) {
	stmt := f.db.MustPrepare(`select name from article_tag group by name order by name`)
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if err == storage.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	items := make([]*model.TagListViewItem, 0)
	var name string
	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, model.NewTagListViewItem(name))
	}
	return items, nil
}
