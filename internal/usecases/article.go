package usecases

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type ArticleInterface interface {
    AccessArticle(articleId model.ArticleId, userId *model.UserId) (model.Article, error)
}
