package usecases

import (
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/domain/repository"
)

type Interface interface {
	AuthInterface
	ArticleInterface
	SearchInterface
	HistoryInterface
	ArticleSubscriptionInterface
	SearchSubscriptionInterface
}

type usecasesThroughRepos struct {
	auth                    AuthInterface
	articleRepo             repository.ArticleRepo
	historyRepo             repository.HistoryRepo
	updatesRepo             repository.UpdatesRepo
	articleSubscriptionRepo repository.ArticleSubscriptionRepo
	searchSubscriptionRepo  repository.SearchSubscriptionRepo
}

func NewUsecases(
	auth AuthInterface,
	articleRepo repository.ArticleRepo,
	historyRepo repository.HistoryRepo,
	updatesRepo repository.UpdatesRepo,
	articleSubscriptionRepo repository.ArticleSubscriptionRepo,
	searchSubscriptionRepo repository.SearchSubscriptionRepo) *usecasesThroughRepos {
	return &usecasesThroughRepos{
		auth:                    auth,
		articleRepo:             articleRepo,
		historyRepo:             historyRepo,
		updatesRepo:             updatesRepo,
		articleSubscriptionRepo: articleSubscriptionRepo,
		searchSubscriptionRepo:  searchSubscriptionRepo,
	}
}

func (u *usecasesThroughRepos) Register(request AuthRequest) (AuthToken, error) {
	return u.auth.Register(request)
}

func (u *usecasesThroughRepos) Login(request AuthRequest) (AuthToken, error) {
	return u.auth.Login(request)
}

func (u *usecasesThroughRepos) Decode(token AuthToken) (model.UserId, error) {
	return u.auth.Decode(token)
}

func (u *usecasesThroughRepos) AccessArticle(articleId model.ArticleId, userId *model.UserId) (model.Article, error) {
	article, err := u.articleRepo.ArticleById(articleId)
	if err != nil {
		return model.Article{}, err
	}
	if userId != nil {
		if err := u.historyRepo.ArticleAccessOccurred(*userId, articleId); err != nil {
			return model.Article{}, err
		}
	}
	return article, nil
}

func (u *usecasesThroughRepos) Search(query model.SearchQuery, userId *model.UserId) (model.SearchResult, error) {
	result, err := u.articleRepo.Search(query, 100)
	if err != nil {
		return model.SearchResult{}, err
	}
	if userId != nil {
		if err := u.historyRepo.SearchAccessOccurred(*userId, query.Query); err != nil {
			return model.SearchResult{}, err
		}
	}
	return result, err
}

func (u *usecasesThroughRepos) GetSearchHistory(id model.UserId) (model.UserSearchHistory, error) {
	queries, err := u.historyRepo.GetSearchHistory(id)
	if err != nil {
		return model.UserSearchHistory{}, err
	}
	return model.UserSearchHistory{
		UserId:  id,
		Queries: queries,
	}, nil
}

func (u *usecasesThroughRepos) ClearSearchHistory(id model.UserId) error {
	return u.historyRepo.ClearSearchHistory(id)
}

func (u *usecasesThroughRepos) GetArticleHistory(id model.UserId) (model.UserArticleHistory, error) {
	articles, err := u.historyRepo.GetArticleHistory(id)
	if err != nil {
		return model.UserArticleHistory{}, err
	}
	metas := make([]model.ArticleMeta, 0, len(articles))
	for _, aid := range articles {
		if meta, err := u.articleRepo.ArticleMetaById(aid); err == nil {
			metas = append(metas, meta)
		} else {
			return model.UserArticleHistory{}, err
		}
	}
	return model.UserArticleHistory{
		UserId:   id,
		Articles: metas,
	}, nil
}

func (u *usecasesThroughRepos) ClearArticleHistory(id model.UserId) error {
	return u.historyRepo.ClearArticleHistory(id)
}

func (u *usecasesThroughRepos) GetArticleLastAccess(userId model.UserId, articleId model.ArticleId) (model.UserArticleAccess, error) {
	ts, err := u.historyRepo.GetArticleLastAccessTimestamp(userId, articleId)
	if err != nil {
		return model.UserArticleAccess{}, err
	}
	return model.UserArticleAccess{
		UserId:    userId,
		ArticleId: articleId,
		Timestamp: ts,
	}, nil
}

func (u *usecasesThroughRepos) GetSearchLastAccess(userId model.UserId, query string) (model.UserSearchAccess, error) {
	ts, err := u.historyRepo.GetSearchLastAccessTimestamp(userId, query)
	if err != nil {
		return model.UserSearchAccess{}, err
	}
	return model.UserSearchAccess{
		UserId:    userId,
		Query:     query,
		Timestamp: ts,
	}, nil
}

func (u *usecasesThroughRepos) SubscribeForArticle(userId model.UserId, articleId model.ArticleId) (model.UserArticleSubscription, error) {
	err := u.articleSubscriptionRepo.SubscribeForArticle(userId, articleId)
	if err != nil {
		return model.UserArticleSubscription{}, err
	}
	return model.UserArticleSubscription{
		UserId:    userId,
		ArticleId: articleId,
	}, nil
}

func (u *usecasesThroughRepos) UnsubscribeFromArticle(userId model.UserId, articleId model.ArticleId) error {
	return u.articleSubscriptionRepo.UnsubscribeFromArticle(userId, articleId)
}

func (u *usecasesThroughRepos) CheckArticleSubscription(userId model.UserId, articleId model.ArticleId) (*model.UserArticleSubscription, error) {
	if s, err := u.articleSubscriptionRepo.IsSubscribedForArticle(userId, articleId); err != nil {
		return nil, err
	} else {
		if s {
			return &model.UserArticleSubscription{
				UserId:    userId,
				ArticleId: articleId,
			}, nil
		} else {
			return nil, nil
		}
	}
}

func (u *usecasesThroughRepos) GetArticleSubscriptions(userId model.UserId) ([]model.UserArticleSubscription, error) {
	subs, err := u.articleSubscriptionRepo.GetArticleSubscriptions(userId)
	if err != nil {
		return nil, err
	}
	result := make([]model.UserArticleSubscription, len(subs))
	for i := range subs {
		result[i] = model.UserArticleSubscription{
			UserId:    userId,
			ArticleId: subs[i],
		}
	}
	return result, nil
}

func (u *usecasesThroughRepos) GetArticleUpdates(userId model.UserId) ([]model.ArticleMeta, error) {
	return u.updatesRepo.GetArticleSubscriptionsUpdates(userId)
}

func (u *usecasesThroughRepos) SubscribeForSearch(userId model.UserId, query string) (model.UserSearchSubscription, error) {
	err := u.searchSubscriptionRepo.SubscribeForSearch(userId, query)
	if err != nil {
		return model.UserSearchSubscription{}, err
	}
	return model.UserSearchSubscription{
		UserId: userId,
		Query:  query,
	}, nil
}

func (u *usecasesThroughRepos) UnsubscribeFromSearch(userId model.UserId, query string) error {
	return u.searchSubscriptionRepo.UnsubscribeFromSearch(userId, query)
}

func (u *usecasesThroughRepos) CheckSearchSubscription(userId model.UserId, query string) (*model.UserSearchSubscription, error) {
	if s, err := u.searchSubscriptionRepo.IsSubscribedForSearch(userId, query); err != nil {
		return nil, err
	} else {
		if s {
			return &model.UserSearchSubscription{
				UserId: userId,
				Query:  query,
			}, nil
		} else {
			return nil, nil
		}
	}
}

func (u *usecasesThroughRepos) GetSearchSubscriptions(userId model.UserId) ([]model.UserSearchSubscription, error) {
	qs, err := u.searchSubscriptionRepo.GetSearchSubscriptions(userId)
	if err != nil {
		return nil, err
	}
	result := make([]model.UserSearchSubscription, len(qs))
	for i := range qs {
		result[i] = model.UserSearchSubscription{
			UserId: userId,
			Query:  qs[i],
		}
	}
	return result, nil
}

func (u *usecasesThroughRepos) GetSearchUpdates(userId model.UserId) ([]string, error) {
	return u.updatesRepo.GetSearchSubscriptionsUpdates(userId)
}
