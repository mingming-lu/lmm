package model

import (
	"regexp"

	"lmm/api/domain/model"
	"lmm/api/service/article/domain"
	"lmm/api/strings"
)

var (
	patternArticleTitle   = regexp.MustCompile("[\u4e00-\u9fa5ぁ-んァ-ンa-zA-Z0-9-_ ]")
	articleTitleMaxLength = 140
)

// Text is the text content model of a article
type Text struct {
	model.ValueObject
	title string
	body  string
}

// NewText create a new text model
func NewText(title, body string) (*Text, error) {
	text := Text{}

	if err := text.SetTitle(title); err != nil {
		return nil, err
	}

	if err := text.SetBody(body); err != nil {
		return nil, err
	}

	return &text, nil
}

// Title returns the title of the article text
func (t *Text) Title() string {
	return t.title
}

// Body returns the body of the article text
func (t *Text) Body() string {
	return t.body
}

// Equals compares two article texts
func (t *Text) Equals(other *Text) bool {
	return (t.Title() == other.Title()) && (t.Body() == other.Body())
}

// SetTitle sets title to the text
func (t *Text) SetTitle(title string) error {
	newTitle := strings.TrimSpace(title)
	if newTitle == "" {
		return domain.ErrEmptyArticleTitle
	}
	if len(newTitle) > articleTitleMaxLength {
		return domain.ErrArticleTitleTooLong
	}
	if !patternArticleTitle.MatchString(newTitle) {
		return domain.ErrInvalidArticleTitle
	}
	t.title = newTitle
	return nil
}

// SetBody sets body to the text
func (t *Text) SetBody(newBody string) error {
	t.body = newBody
	return nil
}
