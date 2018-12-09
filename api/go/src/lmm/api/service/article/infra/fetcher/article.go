package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/finder"
	"lmm/api/service/article/domain/model"
	"lmm/api/storage/db"
)

// ArticleFetcher implements domain.model.finder.ArticleFinder
type ArticleFetcher struct {
	db db.DB
}

// NewArticleFetcher creates new ArticleFetcher
func NewArticleFetcher(db db.DB) *ArticleFetcher {
	return &ArticleFetcher{db: db}
}

// ListByPage implementation
func (f *ArticleFetcher) ListByPage(c context.Context, count, page uint, filter finder.ArticleFilter) (*model.ArticleListView, error) {
	countArticlesSQL := `select count(a.id) from article a`
	if filter.Tag != nil {
		countArticlesSQL += ` inner join article_tag t on a.id = t.article where t.name = ?`
	}

	countArticles := f.db.Prepare(c, countArticlesSQL)
	defer countArticles.Close()

	fetchArticlesSQL := `select a.uid, a.title, a.created_at from article a`
	if filter.Tag != nil {
		fetchArticlesSQL += ` inner join article_tag t on a.id = t.article where t.name = ?`
	}
	fetchArticlesSQL += ` order by created_at desc limit ? offset ?`

	fetchArticles := f.db.Prepare(c, fetchArticlesSQL)
	defer fetchArticles.Close()

	args := make([]interface{}, 0, 3)
	if filter.Tag != nil {
		args = append(args, *filter.Tag)
	}

	var total uint
	if err := countArticles.QueryRow(c, args...).Scan(&total); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql: %s, args: %#v", countArticlesSQL, args))
	}

	args = append(args, count+1, (page-1)*count)
	rows, err := fetchArticles.Query(c, args...)
	if err != nil {
		if err == db.ErrNoRows {
			return model.NewArticleListView(nil, 0, 0, 0, false), nil
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

	return model.NewArticleListView(articles, page, count, total, hasNextPage), nil
}

// FindByID implementation
func (f *ArticleFetcher) FindByID(c context.Context, id *model.ArticleID) (*model.ArticleView, error) {
	selectArticle := f.db.Prepare(c, `
		select id, uid, title, body, created_at, updated_at from article where uid = ?
	`)
	defer selectArticle.Close()

	selectTags := f.db.Prepare(c, `
		select sort, name from article_tag where article = ?
	`)
	defer selectTags.Close()

	var (
		linkedID        uint
		rawArticleID    string
		articleTitle    string
		articleBody     string
		articlePostAt   time.Time
		articleEditedAt time.Time
	)

	err := selectArticle.QueryRow(c, id.String()).Scan(&linkedID, &rawArticleID, &articleTitle, &articleBody, &articlePostAt, &articleEditedAt)
	if err != nil {
		if err == db.ErrNoRows {
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

	rows, err := selectTags.Query(c, linkedID)
	if err != nil && err != db.ErrNoRows {
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
func (f *ArticleFetcher) ListAllTags(c context.Context) (model.TagListView, error) {
	stmt := f.db.Prepare(c, `
		select name from article_tag group by name order by name
	`)
	defer stmt.Close()

	rows, err := stmt.Query(c)
	if err != nil {
		if err == db.ErrNoRows {
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
