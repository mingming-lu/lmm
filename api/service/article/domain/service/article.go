package service

import (
	"context"

	"lmm/api/service/article/domain/model"
	"lmm/api/service/article/domain/repository"
)

// ArticleService is a domain service to deal with article
type ArticleService struct {
	articleRepo repository.ArticleRepository
}

// NewArticleService constructs a new ArticleService
func NewArticleService(articleRepo repository.ArticleRepository) *ArticleService {
	return &ArticleService{articleRepo: articleRepo}
}

// NewArticleToPost build a new article that author wants to post
func (s *ArticleService) NewArticleToPost(c context.Context, author *model.Author, title, body string, tagNames []string) (*model.Article, error) {
	articleID, err := model.NewArticleID(s.articleRepo.NextID(c))
	if err != nil {
		return nil, err
	}

	text, err := model.NewText(title, body)
	if err != nil {
		return nil, err
	}

	tags := make([]*model.Tag, len(tagNames), len(tagNames))
	for i, tagName := range tagNames {
		tag, err := model.NewTag(articleID, uint(i+1), tagName)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	content, err := model.NewContent(text, tags)
	if err != nil {
		return nil, err
	}

	return model.NewArticle(articleID, author, content), nil
}