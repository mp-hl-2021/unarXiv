package dummyUsecases

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "github.com/mp-hl-2021/unarXiv/internal/usecases"
)

var dummyToken usecases.AuthToken = "jwt"
var dummyArticle = model.ArticleMeta{
    Id:                  "dummy",
    Title:               "bunny",
    Authors:             []string{"ya"},
    Abstract:            "abstract",
    LastUpdateTimestamp: 0,
}
var dummyArticleSubscription = model.UserArticleSubscription{
    UserId:    "0",
    ArticleId: "dummy",
}
var dummySearchSubscription = model.UserSearchSubscription{
    UserId: "0",
    Query:  "dummy",
}

type DummyUsecases struct{}

func (d *DummyUsecases) Register(request usecases.AuthRequest) (usecases.AuthToken, error) {
    return dummyToken, nil
}

func (d *DummyUsecases) Login(request usecases.AuthRequest) (usecases.AuthToken, error) {
    return dummyToken, nil
}

func (d *DummyUsecases) Decode(token usecases.AuthToken) (model.UserId, error) {
    return "0", nil
}

func (d *DummyUsecases) AccessArticle(articleId model.ArticleId, userId *model.UserId) (model.Article, error) {
    return model.Article{
        ArticleMeta: dummyArticle,
    }, nil
}

func (d *DummyUsecases) Search(query model.SearchQuery, userId *model.UserId) (model.SearchResult, error) {
    return model.SearchResult{
        TotalMatchesCount: 3,
        Articles: []model.ArticleMeta{dummyArticle},
    }, nil
}

func (d *DummyUsecases) GetSearchHistory(id model.UserId) (model.UserSearchHistory, error) {
    return model.UserSearchHistory{
        UserId:  "0",
        Queries: []string{"dummy"},
    }, nil
}

func (d *DummyUsecases) ClearSearchHistory(id model.UserId) error {
    return nil
}

func (d *DummyUsecases) GetArticleHistory(id model.UserId) (model.UserArticleHistory, error) {
    return model.UserArticleHistory{
        UserId:   "0",
        Articles: []model.ArticleMeta{dummyArticle},
    }, nil
}

func (d *DummyUsecases) ClearArticleHistory(id model.UserId) error {
    return nil
}

func (d *DummyUsecases) GetArticleLastAccess(userId model.UserId, articleId model.ArticleId) (model.UserArticleAccess, error) {
	return model.UserArticleAccess{
		UserId:    "0",
		ArticleId: "dummy",
		Timestamp: 293,
	}, nil
}

func (d *DummyUsecases) GetSearchLastAccess(userId model.UserId, query string) (model.UserSearchAccess, error) {
    return model.UserSearchAccess{
        UserId:    "0",
        Query:     "dummy",
        Timestamp: 2394,
    }, nil
}

func (d *DummyUsecases) SubscribeForArticle(userId model.UserId, articleId model.ArticleId) (model.UserArticleSubscription, error) {
    return dummyArticleSubscription, nil
}

func (d *DummyUsecases) UnsubscribeFromArticle(userId model.UserId, articleId model.ArticleId) error {
    return nil
}

func (d *DummyUsecases) CheckArticleSubscription(userId model.UserId, articleId model.ArticleId) (*model.UserArticleSubscription, error) {
    return &dummyArticleSubscription, nil
}

func (d *DummyUsecases) GetArticleSubscriptions(userId model.UserId) ([]model.UserArticleSubscription, error) {
    return []model.UserArticleSubscription{dummyArticleSubscription}, nil
}

func (d *DummyUsecases) GetArticleUpdates(userId model.UserId) ([]model.ArticleMeta, error) {
    return []model.ArticleMeta{dummyArticle}, nil
}

func (d *DummyUsecases) SubscribeForSearch(userId model.UserId, query string) (model.UserSearchSubscription, error) {
    return dummySearchSubscription, nil
}

func (d *DummyUsecases) UnsubscribeFromSearch(userId model.UserId, query string) error {
    return nil
}

func (d *DummyUsecases) CheckSearchSubscription(userId model.UserId, query string) (*model.UserSearchSubscription, error) {
    return &dummySearchSubscription, nil
}

func (d *DummyUsecases) GetSearchSubscriptions(userId model.UserId) ([]model.UserSearchSubscription, error) {
    return []model.UserSearchSubscription{dummySearchSubscription}, nil
}

func (d *DummyUsecases) GetSearchUpdates(userId model.UserId) ([]string, error) {
    return []string{"dummy"}, nil
}
