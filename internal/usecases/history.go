package usecases

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type HistoryInterface interface {
    GetSearchHistory(id model.UserId) (model.UserSearchHistory, error)
    ClearSearchHistory(id model.UserId) error

    GetArticleHistory(id model.UserId) (model.UserArticleHistory, error)
    ClearArticleHistory(id model.UserId) error

    GetArticleLastAccess(userId model.UserId, articleId model.ArticleId) (*model.UserArticleAccess, error)
    GetSearchLastAccess(userId model.UserId, query string) (*model.UserSearchAccess, error)
}
