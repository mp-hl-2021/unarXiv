package core

import "github.com/mp-hl-2021/unarXiv/api"

var dummyUint32 uint32 = 0
var dummyUint64 uint64 = 0

var dummyAuthenticationRequest = api.AuthenticationRequest{
	Login:    "login",
	Password: "password",
}

var dummyAuthenticationData = api.AuthenticationData{
	Jwt: "jwt",
}

var dummyArticleUserMetaInfo = api.ArticleUserMetaInfo{
	IsSubscribed: true,
}

var dummyArticleMetaInfo = api.ArticleMetaInfo{
	Title:        "title",
	Authors:      "authors",
	ArticleId:    "articleId",
	Abstract:     "abstract",
	UserMetaInfo: &dummyArticleUserMetaInfo,
}

var dummySearchQueryRequest = api.SearchQueryRequest{
	Query:    "query",
	Offset:   &dummyUint32,
	AuthData: &dummyAuthenticationData,
}

var dummySearchQueryResponseItem = api.SearchQueryResponseItem{
	Article: dummyArticleMetaInfo,
}

var dummySearchQueryUserMetaInfo = api.SearchQueryUserMetaInfo{
	IsSubscribed: true,
}

var dummySearchQueryResponse = api.SearchQueryResponse{
	TotalMatchesCount: 1,
	Matches:           []api.SearchQueryResponseItem{dummySearchQueryResponseItem},
	UserMetaInfo:      &dummySearchQueryUserMetaInfo,
}

var dummyGetArticleSubscriptionStatusRequest = api.GetArticleSubscriptionStatusRequest{
	AuthData:  dummyAuthenticationData,
	ArticleId: "articleId",
}

var dummySetSearchQuerySubscriptionStatusRequest = api.SetSearchQuerySubscriptionStatusRequest{
	AuthData:  dummyAuthenticationData,
	Query:     "query",
	Subscribe: true,
}

var dummyGetSearchQuerySubscriptionStatusRequest = api.GetSearchQuerySubscriptionStatusRequest{
	AuthData: dummyAuthenticationData,
	Query:    "query",
}

var dummyGetArticlesUpdatesResponse = api.GetArticlesUpdatesResponse{
	UpdatedArticles: []api.ArticleMetaInfo{dummyArticleMetaInfo},
}

var dummyAccessArticleRequest = api.AccessArticleRequest{
	AuthData:  dummyAuthenticationData,
	ArticleId: "articleId",
}

var dummyAccessArticleResponse = api.AccessArticleResponse{
	Article: dummyArticleMetaInfo,
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
	UserMetaInfo: dummySearchQueryUserMetaInfo,
}

var dummyGetAccessedArticlesHistoryResponse = api.GetAccessedArticlesHistoryResponse{
	AccessedArticlesHistory: []api.ArticleMetaInfo{dummyArticleMetaInfo},
}

type DummyUnarXivAPI struct{}

func (DummyUnarXivAPI) Register(*api.AuthenticationRequest) (*api.AuthenticationData, error) {
	return &dummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Login(*api.AuthenticationRequest) (*api.AuthenticationData, error) {
	return &dummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Search(*api.SearchQueryRequest) (*api.SearchQueryResponse, error) {
	return &dummySearchQueryResponse, nil
}

func (DummyUnarXivAPI) SetArticleSubscriptionStatus(*api.SetArticleSubscriptionStatusRequest) (*api.ArticleUserMetaInfo, error) {
	return &dummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetArticleSubscriptionStatus(*api.GetArticleSubscriptionStatusRequest) (*api.ArticleUserMetaInfo, error) {
	return &dummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) SetSearchQuerySubscriptionStatus(*api.SetSearchQuerySubscriptionStatusRequest) (*api.SearchQueryUserMetaInfo, error) {
	return &dummySearchQueryUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetSearchQuerySubscriptionStatus(*api.GetSearchQuerySubscriptionStatusRequest) (*api.SearchQueryUserMetaInfo, error) {
	return &dummySearchQueryUserMetaInfo, nil
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
