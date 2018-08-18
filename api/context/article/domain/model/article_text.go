package model

import (
	"regexp"

	"lmm/api/context/article/domain"
	"lmm/api/domain/model"
	"lmm/api/utils/strings"
)

var (
	patternArticleTitle = regexp.MustCompile("^[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]$")
)

type ArticleText struct {
	model.ValueObject
	title string
	body  string
}

func NewArticleText(title, body string) (ArticleText, error) {
	text := ArticleText{}
	if err := text.SetTitle(title); err != nil {
		return text, err
	}
	text.SetBody(body)
	return text, nil
}

func (t ArticleText) Title() string {
	return t.title
}

func (t ArticleText) Body() string {
	return t.body
}

func (t ArticleText) Equals(o ArticleText) bool {
	return (t.Title() == o.Title()) && (t.Body() == o.Body())
}

func (t ArticleText) SetTitle(title string) error {
	newTitle := strings.TrimSpace(title)
	if newTitle == "" {
		return domain.ErrEmptyArticleTitle
	}
	if len(newTitle) > 30 {
		return domain.ErrArticleTitleTooLong
	}
	if !patternArticleTitle.MatchString(newTitle) {
		return domain.ErrInvalidArticleTitle
	}
	t.title = newTitle
	return nil
}

func (t ArticleText) SetBody(newBody string) {
	t.body = newBody
}
