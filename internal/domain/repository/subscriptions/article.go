package subscriptions

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type ArticleSubscriptionRepo interface {
    GetArticleSubscriptions(id model.UserId) ([]model.ArticleId, error)
    SubscribeForArticle(id model.UserId, articleId model.ArticleId) error
    UnsubscribeFromArticle(id model.UserId, articleId model.ArticleId) error
    IsSubscribedForArticle(id model.UserId, articleId model.ArticleId) (bool, error)
}
