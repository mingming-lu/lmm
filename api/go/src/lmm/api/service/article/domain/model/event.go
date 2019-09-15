package model

import "context"

type ArticleEventPublisher interface {
	NotifyArticlePublished(c context.Context, id *ArticleID) error
}
