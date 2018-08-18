package model

import "lmm/api/domain/model"

type ArticleWriter struct {
	model.ValueObject
	name string
}

func NewArticleWriter(name string) ArticleWriter {
	return ArticleWriter{name: name}
}

func (w ArticleWriter) Name() string {
	return w.name
}
