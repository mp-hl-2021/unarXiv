package repository

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type UpdatesRepo interface {
    GetArticleSubscriptionsUpdates(id model.UserId) ([]model.ArticleMeta, error)
    GetSearchSubscriptionsUpdates(id model.UserId) ([]string, error)
}
