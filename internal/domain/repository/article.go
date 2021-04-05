package repository

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
)

type ArticleRepo interface {
    ArticleMetaById(id model.ArticleId) (model.ArticleMeta, error)
    ArticleById(id model.ArticleId) (model.Article, error)

    UpdateArticle(article model.Article) error

    Search(query model.SearchQuery, limit uint32) (model.SearchResult, error)
}
