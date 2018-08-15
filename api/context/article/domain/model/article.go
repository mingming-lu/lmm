package model

import (
	"lmm/api/context/article/domain"
	"lmm/api/utils/rand"
	"lmm/api/utils/strings"

	"regexp"
)

var (
	patternArticleID    = regexp.MustCompile("[0-9a-Z]{6}")
	patternArticleTitle = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]$")
)

type ArticleID struct {
	id string
}

var emptyArticleID ArticleID = ArticleID{id: ""}

func NewArticleID(s string) (ArticleID, error) {
	if !patternArticleID.MatchString(s) {
		return emptyArticleID, domain.ErrInvalidArticleID
	}
	return ArticleID{id: s}, nil
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
