package usecases

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type ArticleUserRelationsInterface interface {
	SubscribeForArticle(userId model.UserId, articleId model.ArticleId) (model.UserArticleSubscription, error)
	UnsubscribeFromArticle(userId model.UserId, articleId model.ArticleId) error
	CheckArticleSubscription(userId model.UserId, articleId model.ArticleId) (*model.UserArticleSubscription, error)

	GetArticleSubscriptions(userId model.UserId) ([]model.UserArticleSubscription, error)

	GetArticleUpdates(userId model.UserId) ([]model.ArticleMeta, error)

	GetArticleHistory(id model.UserId) (model.UserArticleHistory, error)
	ClearArticleHistory(id model.UserId) error
	GetArticleLastAccess(userId model.UserId, articleId model.ArticleId) (model.UserArticleAccess, error)
}
