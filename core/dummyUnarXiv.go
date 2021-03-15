package core

import "github.com/mp-hl-2021/unarXiv/api"

var dummyUint32 uint32 = 0
var dummyUint64 uint64 = 0

var DummyAuthenticationData = api.AuthenticationData{
    Jwt: "jwt",
}

var DummyArticleUserMetaInfo = api.ArticleUserMetaInfo{
    IsSubscribed: true,
}

var DummyArticleMetaInfo = api.ArticleMetaInfo{
    Title:        "title",
    Authors:      "authors",
    ArticleId:    "articleId",
    Abstract:     "abstract",
    UserMetaInfo: &DummyArticleUserMetaInfo,
}

var dummySearchQueryResponseItem = api.SearchQueryResponseItem{
    Article: DummyArticleMetaInfo,
}

var DummySearchQueryUserMetaInfo = api.SearchQueryUserMetaInfo{
    IsSubscribed: true,
}

var DummySearchQueryResponse = api.SearchQueryResponse{
    TotalMatchesCount: 1,
    Matches:           []api.SearchQueryResponseItem{dummySearchQueryResponseItem},
    UserMetaInfo:      &DummySearchQueryUserMetaInfo,
}

var dummyGetArticlesUpdatesResponse = api.GetArticlesUpdatesResponse{
    UpdatedArticles: []api.ArticleMetaInfo{DummyArticleMetaInfo},
}

var dummyAccessArticleResponse = api.AccessArticleResponse{
    Article: DummyArticleMetaInfo,
}

var dummyGetSearchQueriesUpdatesResponse = api.GetSearchQueriesUpdatesResponse{
    SearchQueriesUpdates: []api.GetSearchQueriesUpdatesResponseItem{dummyGetSearchQueriesUpdatesResponseItem},
}

var dummyGetSearchQueriesUpdatesResponseItem = api.GetSearchQueriesUpdatesResponseItem{
    Query:                   "query",
    LastNewArticleTimestamp: &dummyUint64,
    NewArticlesCount:        &dummyUint32,
}

var dummyGetSearchHistoryResponse = api.GetSearchHistoryResponse{
    SearchHistory: []api.GetSearchHistoryResponseItem{dummyGetSearchHistoryResponseItem},
}

var dummyGetSearchHistoryResponseItem = api.GetSearchHistoryResponseItem{
    Query:        "query",
    UserMetaInfo: DummySearchQueryUserMetaInfo,
}

var dummyGetAccessedArticlesHistoryResponse = api.GetAccessedArticlesHistoryResponse{
    AccessedArticlesHistory: []api.ArticleMetaInfo{DummyArticleMetaInfo},
}

type DummyUnarXivAPI struct{}

func (DummyUnarXivAPI) Register(*api.AuthenticationRequest) (*api.AuthenticationData, error) {
    return &DummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Login(*api.AuthenticationRequest) (*api.AuthenticationData, error) {
    return &DummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Search(*api.SearchQueryRequest) (*api.SearchQueryResponse, error) {
    return &DummySearchQueryResponse, nil
}

func (DummyUnarXivAPI) SetArticleSubscriptionStatus(*api.SetArticleSubscriptionStatusRequest) (*api.ArticleUserMetaInfo, error) {
    return &DummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetArticleSubscriptionStatus(*api.GetArticleSubscriptionStatusRequest) (*api.ArticleUserMetaInfo, error) {
    return &DummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) SetSearchQuerySubscriptionStatus(*api.SetSearchQuerySubscriptionStatusRequest) (*api.SearchQueryUserMetaInfo, error) {
    return &DummySearchQueryUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetSearchQuerySubscriptionStatus(*api.GetSearchQuerySubscriptionStatusRequest) (*api.SearchQueryUserMetaInfo, error) {
    return &DummySearchQueryUserMetaInfo, nil
}

func (DummyUnarXivAPI) GetArticlesUpdates(*api.AuthenticationData) (*api.GetArticlesUpdatesResponse, error) {
    return &dummyGetArticlesUpdatesResponse, nil
}
func (DummyUnarXivAPI) AccessArticle(*api.AccessArticleRequest) (*api.AccessArticleResponse, error) {
    return &dummyAccessArticleResponse, nil
}

func (DummyUnarXivAPI) GetSearchQueriesUpdates(*api.AuthenticationData) (*api.GetSearchQueriesUpdatesResponse, error) {
    return &dummyGetSearchQueriesUpdatesResponse, nil
}

func (DummyUnarXivAPI) GetSearchHistory(*api.AuthenticationData) (*api.GetSearchHistoryResponse, error) {
    return &dummyGetSearchHistoryResponse, nil
}
func (DummyUnarXivAPI) GetArticlesHistory(*api.AuthenticationData) (*api.GetAccessedArticlesHistoryResponse, error) {
    return &dummyGetAccessedArticlesHistoryResponse, nil
}
