package service

import (
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/repository"
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
func (s *ArticleService) NewArticleToPost(author *model.Author, title, body string, tagNames []string) (*model.Article, error) {
	articleID, err := model.NewArticleID(s.articleRepo.NextID())
	if err != nil {
		return nil, err
	}

	text, err := model.NewText(title, body)
	if err != nil {
		return nil, err
	}

	tags := make([]*model.Tag, len(tagNames), len(tagNames))

	for i, tagName := range tagNames {
		tag, err := model.NewTag(articleID, tagName)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	return model.NewArticle(articleID, text, author, tags), nil
}
