package core

import "github.com/mp-hl-2021/unarXiv/usecases"

var dummyUint32 uint32 = 0
var dummyUint64 uint64 = 0

var DummyAuthenticationData = usecases.AuthenticationData{
    Jwt: "jwt",
}

var DummyArticleUserMetaInfo = usecases.ArticleUserMetaInfo{
    IsSubscribed: true,
}

var DummyArticleMetaInfo = usecases.ArticleMetaInfo{
    Title:        "title",
    Authors:      "authors",
    ArticleId:    "articleId",
    Abstract:     "abstract",
    UserMetaInfo: &DummyArticleUserMetaInfo,
}

var dummySearchQueryResponseItem = usecases.SearchQueryResponseItem{
    Article: DummyArticleMetaInfo,
}

var DummySearchQueryUserMetaInfo = usecases.SearchQueryUserMetaInfo{
    IsSubscribed: true,
}

var DummySearchQueryResponse = usecases.SearchQueryResponse{
    TotalMatchesCount: 1,
    Matches:           []usecases.SearchQueryResponseItem{dummySearchQueryResponseItem},
    UserMetaInfo:      &DummySearchQueryUserMetaInfo,
}

var dummyGetArticlesUpdatesResponse = usecases.GetArticlesUpdatesResponse{
    UpdatedArticles: []usecases.ArticleMetaInfo{DummyArticleMetaInfo},
}

var dummyAccessArticleResponse = usecases.AccessArticleResponse{
    Article: DummyArticleMetaInfo,
}

var dummyGetSearchQueriesUpdatesResponse = usecases.GetSearchQueriesUpdatesResponse{
    SearchQueriesUpdates: []usecases.GetSearchQueriesUpdatesResponseItem{dummyGetSearchQueriesUpdatesResponseItem},
}

var dummyGetSearchQueriesUpdatesResponseItem = usecases.GetSearchQueriesUpdatesResponseItem{
    Query:                   "query",
    LastNewArticleTimestamp: &dummyUint64,
    NewArticlesCount:        &dummyUint32,
}

var dummyGetSearchHistoryResponse = usecases.GetSearchHistoryResponse{
    SearchHistory: []usecases.GetSearchHistoryResponseItem{dummyGetSearchHistoryResponseItem},
}

var dummyGetSearchHistoryResponseItem = usecases.GetSearchHistoryResponseItem{
    Query:        "query",
    UserMetaInfo: DummySearchQueryUserMetaInfo,
}

var dummyGetAccessedArticlesHistoryResponse = usecases.GetAccessedArticlesHistoryResponse{
    AccessedArticlesHistory: []usecases.ArticleMetaInfo{DummyArticleMetaInfo},
}

type DummyUnarXivAPI struct{}

func (DummyUnarXivAPI) Register(*usecases.AuthenticationRequest) (*usecases.AuthenticationData, error) {
    return &DummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Login(*usecases.AuthenticationRequest) (*usecases.AuthenticationData, error) {
    return &DummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Search(*usecases.SearchQueryRequest) (*usecases.SearchQueryResponse, error) {
    return &DummySearchQueryResponse, nil
}

func (DummyUnarXivAPI) SetArticleSubscriptionStatus(*usecases.SetArticleSubscriptionStatusRequest) (*usecases.ArticleUserMetaInfo, error) {
    return &DummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetArticleSubscriptionStatus(*usecases.GetArticleSubscriptionStatusRequest) (*usecases.ArticleUserMetaInfo, error) {
    return &DummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) SetSearchQuerySubscriptionStatus(*usecases.SetSearchQuerySubscriptionStatusRequest) (*usecases.SearchQueryUserMetaInfo, error) {
    return &DummySearchQueryUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetSearchQuerySubscriptionStatus(*usecases.GetSearchQuerySubscriptionStatusRequest) (*usecases.SearchQueryUserMetaInfo, error) {
    return &DummySearchQueryUserMetaInfo, nil
}

func (DummyUnarXivAPI) GetArticlesUpdates(*usecases.AuthenticationData) (*usecases.GetArticlesUpdatesResponse, error) {
    return &dummyGetArticlesUpdatesResponse, nil
}
func (DummyUnarXivAPI) AccessArticle(*usecases.AccessArticleRequest) (*usecases.AccessArticleResponse, error) {
    return &dummyAccessArticleResponse, nil
}

func (DummyUnarXivAPI) GetSearchQueriesUpdates(*usecases.AuthenticationData) (*usecases.GetSearchQueriesUpdatesResponse, error) {
    return &dummyGetSearchQueriesUpdatesResponse, nil
}

func (DummyUnarXivAPI) GetSearchHistory(*usecases.AuthenticationData) (*usecases.GetSearchHistoryResponse, error) {
    return &dummyGetSearchHistoryResponse, nil
}
func (DummyUnarXivAPI) GetArticlesHistory(*usecases.AuthenticationData) (*usecases.GetAccessedArticlesHistoryResponse, error) {
    return &dummyGetAccessedArticlesHistoryResponse, nil
}
