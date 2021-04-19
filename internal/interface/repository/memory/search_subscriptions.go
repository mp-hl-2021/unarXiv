package memory

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain"
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "sync"
)

type SearchSubscriptionRepo struct {
    subscriptions map[model.UserId][]string
    mutex         *sync.Mutex
}

func NewSearchSubscriptionRepo() *SearchSubscriptionRepo {
    return &SearchSubscriptionRepo{
        subscriptions: make(map[model.UserId][]string),
        mutex:         &sync.Mutex{},
    }
}

func (a *SearchSubscriptionRepo) GetSearchSubscriptions(id model.UserId) ([]string, error) {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    if subs, ok := a.subscriptions[id]; ok {
        return subs, nil
    } else {
        return []string{}, nil
    }
}

func (a *SearchSubscriptionRepo) SubscribeForSearch(id model.UserId, query string) error {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    if subs, ok := a.subscriptions[id]; !ok {
        a.subscriptions[id] = []string{query}
    } else {
        for _, aid := range subs {
            if aid == query {
                return domain.AlreadySubscribed
            }
        }
        a.subscriptions[id] = append(subs, query)
    }
    return nil
}

func (a *SearchSubscriptionRepo) UnsubscribeFromSearch(id model.UserId, query string) error {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    if subs, ok := a.subscriptions[id]; ok {
        for i, aid := range subs {
            if aid == query {
                subs[i] = subs[len(subs)-1]
                a.subscriptions[id] = subs[:len(subs)-1]
                return nil
            }
        }
    }
    return domain.NotSubscribed
}

func (a *SearchSubscriptionRepo) IsSubscribedForSearch(id model.UserId, query string) (bool, error) {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    if subs, ok := a.subscriptions[id]; ok {
        for _, aid := range subs {
            if aid == query {
                return true, nil
            }
        }
    }
    return false, nil
}
