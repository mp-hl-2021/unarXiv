package subscriptions

import (
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"sync"
)

type ArticleSubscriptionRepo struct {
	subscriptions map[model.UserId][]model.ArticleId
	mutex         *sync.Mutex
}

func NewArticleSubscriptionRepo() *ArticleSubscriptionRepo {
	return &ArticleSubscriptionRepo{
		subscriptions: make(map[model.UserId][]model.ArticleId),
		mutex:         &sync.Mutex{},
	}
}

func (a *ArticleSubscriptionRepo) GetArticleSubscriptions(id model.UserId) ([]model.ArticleId, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if subs, ok := a.subscriptions[id]; ok {
		return subs, nil
	} else {
		return []model.ArticleId{}, nil
	}
}

func (a *ArticleSubscriptionRepo) SubscribeForArticle(id model.UserId, articleId model.ArticleId) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if subs, ok := a.subscriptions[id]; !ok {
		a.subscriptions[id] = []model.ArticleId{articleId}
	} else {
		for _, aid := range subs {
			if aid == articleId {
				return domain.AlreadySubscribed
			}
		}
		a.subscriptions[id] = append(subs, articleId)
	}
	return nil
}

func (a *ArticleSubscriptionRepo) UnsubscribeFromArticle(id model.UserId, articleId model.ArticleId) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if subs, ok := a.subscriptions[id]; ok {
		for i, aid := range subs {
			if aid == articleId {
				subs[i] = subs[len(subs)-1]
				a.subscriptions[id] = subs[:len(subs)-1]
				return nil
			}
		}
	}
	return domain.NotSubscribed
}

func (a *ArticleSubscriptionRepo) IsSubscribedForArticle(id model.UserId, articleId model.ArticleId) (bool, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if subs, ok := a.subscriptions[id]; ok {
		for _, aid := range subs {
			if aid == articleId {
				return true, nil
			}
		}
	}
	return false, nil
}
