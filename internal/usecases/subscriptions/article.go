package subscriptions

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type ArticleInterface interface {
    SubscribeForArticle(userId model.UserId, articleId model.ArticleId) (model.UserArticleSubscription, error)
    UnsubscribeFromArticle(userId model.UserId, articleId model.ArticleId) error
    CheckArticleSubscription(userId model.UserId, articleId model.ArticleId) (*model.UserArticleSubscription, error)

    GetArticleSubscriptions(userId model.UserId) ([]model.UserArticleSubscription, error)

    GetArticleUpdates(userId model.UserId) ([]model.ArticleMeta, error)
}
