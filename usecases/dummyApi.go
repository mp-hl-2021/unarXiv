package usecases

var dummyUint32 uint32 = 0
var dummyUint64 uint64 = 0

var DummyAuthenticationData = AuthenticationData{
	Jwt: "jwt",
}

var DummyArticleUserMetaInfo = ArticleUserMetaInfo{
	IsSubscribed: true,
}

var DummyArticleMetaInfo = ArticleMetaInfo{
	Title:        "title",
	Authors:      "authors",
	ArticleId:    "articleId",
	Abstract:     "abstract",
	UserMetaInfo: &DummyArticleUserMetaInfo,
}

var dummySearchQueryResponseItem = SearchQueryResponseItem{
	Article: DummyArticleMetaInfo,
}

var DummySearchQueryUserMetaInfo = SearchQueryUserMetaInfo{
	IsSubscribed: true,
}

var DummySearchQueryResponse = SearchQueryResponse{
	TotalMatchesCount: 1,
	Matches:           []SearchQueryResponseItem{dummySearchQueryResponseItem},
	UserMetaInfo:      &DummySearchQueryUserMetaInfo,
}

var dummyGetArticlesUpdatesResponse = GetArticlesUpdatesResponse{
	UpdatedArticles: []ArticleMetaInfo{DummyArticleMetaInfo},
}

var dummyAccessArticleResponse = AccessArticleResponse{
	Article: DummyArticleMetaInfo,
}

var dummyGetSearchQueriesUpdatesResponse = GetSearchQueriesUpdatesResponse{
	SearchQueriesUpdates: []GetSearchQueriesUpdatesResponseItem{dummyGetSearchQueriesUpdatesResponseItem},
}

var dummyGetSearchQueriesUpdatesResponseItem = GetSearchQueriesUpdatesResponseItem{
	Query:                   "query",
	LastNewArticleTimestamp: &dummyUint64,
	NewArticlesCount:        &dummyUint32,
}

var dummyGetSearchHistoryResponse = GetSearchHistoryResponse{
	SearchHistory: []GetSearchHistoryResponseItem{dummyGetSearchHistoryResponseItem},
}

var dummyGetSearchHistoryResponseItem = GetSearchHistoryResponseItem{
	Query:        "query",
	UserMetaInfo: DummySearchQueryUserMetaInfo,
}

var dummyGetAccessedArticlesHistoryResponse = GetAccessedArticlesHistoryResponse{
	AccessedArticlesHistory: []ArticleMetaInfo{DummyArticleMetaInfo},
}

type DummyUnarXivAPI struct{}

func (DummyUnarXivAPI) Register(*AuthenticationRequest) (*AuthenticationData, error) {
	return &DummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Login(*AuthenticationRequest) (*AuthenticationData, error) {
	return &DummyAuthenticationData, nil
}
func (DummyUnarXivAPI) Search(*SearchQueryRequest) (*SearchQueryResponse, error) {
	return &DummySearchQueryResponse, nil
}

func (DummyUnarXivAPI) SetArticleSubscriptionStatus(*SetArticleSubscriptionStatusRequest) (*ArticleUserMetaInfo, error) {
	return &DummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetArticleSubscriptionStatus(*GetArticleSubscriptionStatusRequest) (*ArticleUserMetaInfo, error) {
	return &DummyArticleUserMetaInfo, nil
}
func (DummyUnarXivAPI) SetSearchQuerySubscriptionStatus(*SetSearchQuerySubscriptionStatusRequest) (*SearchQueryUserMetaInfo, error) {
	return &DummySearchQueryUserMetaInfo, nil
}
func (DummyUnarXivAPI) GetSearchQuerySubscriptionStatus(*GetSearchQuerySubscriptionStatusRequest) (*SearchQueryUserMetaInfo, error) {
	return &DummySearchQueryUserMetaInfo, nil
}

func (DummyUnarXivAPI) GetArticlesUpdates(*AuthenticationData) (*GetArticlesUpdatesResponse, error) {
	return &dummyGetArticlesUpdatesResponse, nil
}
func (DummyUnarXivAPI) AccessArticle(*AccessArticleRequest) (*AccessArticleResponse, error) {
	return &dummyAccessArticleResponse, nil
}

func (DummyUnarXivAPI) GetSearchQueriesUpdates(*AuthenticationData) (*GetSearchQueriesUpdatesResponse, error) {
	return &dummyGetSearchQueriesUpdatesResponse, nil
}

func (DummyUnarXivAPI) GetSearchHistory(*AuthenticationData) (*GetSearchHistoryResponse, error) {
	return &dummyGetSearchHistoryResponse, nil
}
func (DummyUnarXivAPI) GetArticlesHistory(*AuthenticationData) (*GetAccessedArticlesHistoryResponse, error) {
	return &dummyGetAccessedArticlesHistoryResponse, nil
}
