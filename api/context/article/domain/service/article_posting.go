package service

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/article/domain/model"
	"lmm/api/context/article/domain/repository"

	"time"
)

type ArticlePostingService struct {
	articleRepo repository.ArticleRepository
}

func NewArticlePostingService(articleRepo repository.ArticleRepository) *ArticlePostingService {
	return &ArticlePostingService{articleRepo: articleRepo}
}

func (s *ArticlePostingService) PostingArticle(
	user *account.User,
	title, text string,
	tagNames []string,
) (*model.Article, error) {
	articleID, err := model.NewArticleID(s.articleRepo.NextID())
	if err != nil {
		return nil, err
	}

	articleText, err := model.NewArticleText(title, text)
	if err != nil {
		return nil, err
	}

	articleWriter := model.NewArticleWriter(user.Name())

	tags := make([]*model.Tag, 0, len(tagNames))
	for _, tagName := range tagNames {
		tag, err := model.NewTag(articleID, tagName)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	now := time.Now()

	return model.NewArticle(articleID, articleText, articleWriter, now, now, tags), nil
}
