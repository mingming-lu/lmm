package model

import (
	"errors"
	"time"

	"lmm/api/context/base/domain/model"
)

var (
	ErrArticleTextNoChange = errors.New("article text no change")
	ErrTagAlreadyAdded     = errors.New("the tag has already been added")
	ErrTagNotAdded         = errors.New("tag not added")
)

type Article struct {
	model.Entity
	id             ArticleID
	text           ArticleText
	writer         ArticleWriter
	postAt         time.Time
	lastModifiedAt time.Time
	tags           []*Tag
}

func NewArticle(
	articleID ArticleID,
	text ArticleText,
	writer ArticleWriter,
	postAt time.Time,
	lastModifiedAt time.Time,
	tags []*Tag,
) *Article {
	return &Article{
		id:             articleID,
		text:           text,
		writer:         writer,
		postAt:         postAt,
		lastModifiedAt: lastModifiedAt,
		tags:           tags,
	}
}

func GenerateArticleID() (ArticleID, error) {
	return NewArticleID(rand.Base62(6))
}

func (id ArticleID) String() string {
	return id.id
}

type Article struct {
	id    ArticleID
	title string
	text  string
}

func NewArticle(articleID ArticleID, title, text string) (*Article, error) {
	title, err := validateArticleTitle(title)
	if err != nil {
		return nil, err
	}
	return &Article{id: articleID, title: title, text: text}, nil
}

func (a *Article) ID() ArticleID {
	return a.id
}

func (a *Article) Title() string {
	return a.title
}

func (a *Article) Text() string {
	return a.text
}

func (a *Article) ChangeTitle(newTitle string) error {
	title, err := validateArticleTitle(newTitle)
	if err != nil {
		return err
	}
	a.title = title
	return nil
}

func (a *Article) ChangeText(newText string) error {
	a.text = newText
	return nil
}

func validateArticleTitle(s string) (string, error) {
	title := strings.TrimSpace(s)
	if title == "" {
		return "", domain.ErrEmptyArticleTitle
	}
	if len(title) > 30 {
		return "", domain.ErrArticleTitleTooLong
	}
	if !patternArticleTitle.MatchString(s) {
		return "", domain.ErrInvalidArticleTitle
	}
	return title, nil
}
