package repository

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type ArticleSubscriptionRepo interface {
    GetArticleSubscriptions(id model.UserId) ([]model.ArticleId, error)
    SubscribeForArticle(id model.UserId, articleId model.ArticleId) error
    UnsubscribeFromArticle(id model.UserId, articleId model.ArticleId) error
    IsSubscribedForArticle(id model.UserId, articleId model.ArticleId) (bool, error)

    ArticleAccessOccurred(userId model.UserId, articleId model.ArticleId) error
    GetArticleLastAccessTimestamp(userId model.UserId, articleId model.ArticleId) (uint64, error)

    GetArticleHistory(userId model.UserId) ([]model.ArticleId, error)
    ClearArticleHistory(userId model.UserId) error

}
