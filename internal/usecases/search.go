package usecases

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
)

type SearchInterface interface {
    Search(query model.SearchQuery, userId *model.UserId) (model.SearchResult, error)
}
