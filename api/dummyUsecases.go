package api

var dummyUint32 uint32 = 0
var dummyUint64 uint64 = 0

var dummyAuthenticationRequest = AuthenticationRequest{
	"login",
	"password",
}

var dummyAuthenticationData = AuthenticationData{
	"jwt",
}

var dummyArticleUserMetaInfo = ArticleUserMetaInfo{
	true,
	0,
}

var dummyArticleMetaInfo = ArticleMetaInfo{
	"title",
	"authors",
	"articleId",
	"abstract",
	0,
	&dummyArticleUserMetaInfo,
}

var dummySearchQueryRequest = SearchQueryRequest{
	"query",
	&dummyUint32,
	&dummyAuthenticationData,
}

var dummySearchQueryResponseItem = SearchQueryResponseItem{
	dummyArticleMetaInfo,
}

var dummySearchQueryUserMetaInfo = SearchQueryUserMetaInfo{
	true,
	0,
}

var dummySearchQueryResponse = SearchQueryResponse{
	1,
	[]SearchQueryResponseItem{dummySearchQueryResponseItem},
	&dummySearchQueryUserMetaInfo,
}

var dummyGetArticleSubscriptionStatusRequest = GetArticleSubscriptionStatusRequest{
	dummyAuthenticationData,
	"articleId",
}

var dummySetSearchQuerySubscriptionStatusRequest = SetSearchQuerySubscriptionStatusRequest{
	dummyAuthenticationData,
	"query",
	true,
}

var dummyGetSearchQuerySubscriptionStatusRequest = GetSearchQuerySubscriptionStatusRequest{
	dummyAuthenticationData,
	"query",
}

var dummyGetArticlesUpdatesResponse = GetArticlesUpdatesResponse{
	[]ArticleMetaInfo{dummyArticleMetaInfo},
}

var dummyAccessArticleRequest = AccessArticleRequest{
	dummyAuthenticationData,
	"articleId",
}

var dummyAccessArticleResponse = AccessArticleResponse{
	dummyArticleMetaInfo,
}

var dummyGetSearchQueriesUpdatesResponse = GetSearchQueriesUpdatesResponse{
	[]GetSearchQueriesUpdatesResponseItem{dummyGetSearchQueriesUpdatesResponseItem},
}

var dummyGetSearchQueriesUpdatesResponseItem = GetSearchQueriesUpdatesResponseItem{
	"query",
	&dummyUint64,
	0,
	&dummyUint32,
}

var dummyGetSearchHistoryResponse = GetSearchHistoryResponse{
	[]GetSearchHistoryResponseItem{dummyGetSearchHistoryResponseItem},
}

var dummyGetSearchHistoryResponseItem = GetSearchHistoryResponseItem{
	"query",
	dummySearchQueryUserMetaInfo,
}

var dummyGetAccessedArticlesHistoryResponse = GetAccessedArticlesHistoryResponse{
	[]ArticleMetaInfo{dummyArticleMetaInfo},
}

type dummyUseCases struct{}

func (dummyUseCases) Register(*AuthenticationRequest) (*AuthenticationData, error) {
	return &dummyAuthenticationData, nil
}
func (dummyUseCases) Login(*AuthenticationRequest) (*AuthenticationData, error) {
	return &dummyAuthenticationData, nil
}
func (dummyUseCases) Search(*SearchQueryRequest) (*SearchQueryResponse, error) {
	return &dummySearchQueryResponse, nil
}

func (dummyUseCases) SetArticleSubscriptionStatus(*SetArticleSubscriptionStatusRequest) (*ArticleUserMetaInfo, error) {
	return &dummyArticleUserMetaInfo, nil
}
func (dummyUseCases) GetArticleSubscriptionStatus(*GetArticleSubscriptionStatusRequest) (*ArticleUserMetaInfo, error) {
	return &dummyArticleUserMetaInfo, nil
}
func (dummyUseCases) SetSearchQuerySubscriptionStatus(*SetSearchQuerySubscriptionStatusRequest) (*SearchQueryUserMetaInfo, error) {
	return &dummySearchQueryUserMetaInfo, nil
}
func (dummyUseCases) GetSearchQuerySubscriptionStatus(*GetSearchQuerySubscriptionStatusRequest) (*SearchQueryUserMetaInfo, error) {
	return &dummySearchQueryUserMetaInfo, nil
}

func (dummyUseCases) GetArticlesUpdates(*AuthenticationData) (*GetArticlesUpdatesResponse, error) {
	return &dummyGetArticlesUpdatesResponse, nil
}
func (dummyUseCases) AccessArticle(*AccessArticleRequest) (*AccessArticleResponse, error) {
	return &dummyAccessArticleResponse, nil
}

func (dummyUseCases) GetSearchQueriesUpdates(*AuthenticationData) (*GetSearchQueriesUpdatesResponse, error) {
	return &dummyGetSearchQueriesUpdatesResponse, nil
}

func (dummyUseCases) GetSearchHistory(*AuthenticationData) (*GetSearchHistoryResponse, error) {
	return &dummyGetSearchHistoryResponse, nil
}
func (dummyUseCases) GetArticlesHistory(*AuthenticationData) (*GetAccessedArticlesHistoryResponse, error) {
	return &dummyGetAccessedArticlesHistoryResponse, nil
}