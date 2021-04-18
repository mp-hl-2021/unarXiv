package repository

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type HistoryRepo interface {
    ArticleAccessOccurred(userId model.UserId, articleId model.ArticleId) error
    GetArticleLastAccessTimestamp(userId model.UserId, articleId model.ArticleId) (uint64, error)

    GetArticleHistory(userId model.UserId) ([]model.ArticleId, error)
    ClearArticleHistory(userId model.UserId) error

    SearchAccessOccurred(userId model.UserId, query string) error
    GetSearchLastAccessTimestamp(userId model.UserId, query string) (uint64, error)

    GetSearchHistory(userId model.UserId) ([]string, error)
    ClearSearchHistory(userId model.UserId) error
}
