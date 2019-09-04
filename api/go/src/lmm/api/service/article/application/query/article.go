package query

import (
	"gopkg.in/go-playground/validator.v8"
)

type ListArticleQuery struct {
	Page    int    `form:"page,default=1" binding:"min=1"`
	PerPage int    `form:"perPage,default=5" binding:"min=1"`
	Tag     string `form:"tag"`
}

func (q *ListArticleQuery) ValidateErrors(err error) []string {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{err.Error()}
	}

	errStrings := make([]string, len(errors))
	i := 0
	for _, err := range errors {
		switch err.Field {
		case "Page":
			errStrings[i] = "invalid page"
		case "PerPage":
			errStrings[i] = "invalid perPage"
		}
		i++
	}

	return errStrings
}
