package subscriptions

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type SearchSubscriptionRepo interface {
    GetSearchSubscriptions(id model.UserId) ([]string, error)
    SubscribeForSearch(id model.UserId, query string) error
    UnsubscribeFromSearch(id model.UserId, query string) error
    IsSubscribedForSearch(id model.UserId, query string) (bool, error)
}
