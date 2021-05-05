package repository

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type SearchUserRelationsRepo interface {
	GetSearchSubscriptions(id model.UserId) ([]string, error)
	SubscribeForSearch(id model.UserId, query string) error
	UnsubscribeFromSearch(id model.UserId, query string) error
	IsSubscribedForSearch(id model.UserId, query string) (bool, error)

	SearchAccessOccurred(userId model.UserId, query string) error
	GetSearchLastAccessTimestamp(userId model.UserId, query string) (uint64, error)

	GetSearchHistory(userId model.UserId) ([]string, error)
	ClearSearchHistory(userId model.UserId) error
}
