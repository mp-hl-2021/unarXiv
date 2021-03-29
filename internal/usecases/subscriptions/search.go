package subscriptions

import "github.com/mp-hl-2021/unarXiv/internal/domain/model"

type SearchInterface interface {
    SubscribeForSearch(userId model.UserId, query string) (model.UserSearchSubscription, error)
    UnsubscribeFromSearch(userId model.UserId, query string) error
    CheckSearchSubscription(userId model.UserId, query string) (*model.UserSearchSubscription, error)

    GetSearchSubscriptions(userId model.UserId) ([]model.UserSearchSubscription, error)

    GetSearchUpdates(userId model.UserId) ([]string, error)
}
