package api

type UnarXivApi interface {
	Register(*AuthenticationRequest) (*AuthenticationData, error)
	Login(*AuthenticationRequest) (*AuthenticationData, error)

	Search(*SearchQueryRequest) (*SearchQueryResponse, error)

	SetArticleSubscriptionStatus(*SetArticleSubscriptionStatusRequest) (*ArticleUserMetaInfo, error)
	GetArticleSubscriptionStatus(*GetArticleSubscriptionStatusRequest) (*ArticleUserMetaInfo, error)
	SetSearchQuerySubscriptionStatus(*SetSearchQuerySubscriptionStatusRequest) (*SearchQueryUserMetaInfo, error)
	GetSearchQuerySubscriptionStatus(*GetSearchQuerySubscriptionStatusRequest) (*SearchQueryUserMetaInfo, error)

	GetArticlesUpdates(*AuthenticationData) (*GetArticlesUpdatesResponse, error)
	AccessArticle(*AccessArticleRequest) (*AccessArticleResponse, error)

	GetSearchQueriesUpdates(*AuthenticationData) (*GetSearchQueriesUpdatesResponse, error)

	GetSearchHistory(*AuthenticationData) (*GetSearchHistoryResponse, error)
	GetArticlesHistory(*AuthenticationData) (*GetAccessedArticlesHistoryResponse, error)
}

type AuthenticationRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthenticationData struct {
	Jwt string `json:"jwt"`
}

type ArticleMetaInfo struct {
	Title     string `json:"title"`
	Authors   string `json:"authors"`
	ArticleId string `json:"articleId"`
	Abstract  string `json:"abstract"`

	LastUpdateTimestamp uint64 `json:"lastUpdateTimestamp"`

	// this field is present if user is authenticated
	UserMetaInfo *ArticleUserMetaInfo `json:"userMetaInfo,omitempty"`
}

type ArticleUserMetaInfo struct {
	IsSubscribed          bool   `json:"isSubscribed"`
	LastAccessedTimestamp uint64 `json:"lastAccessedTimestamp"`
}

type SearchQueryRequest struct {
	Query    string              `json:"query"`
	Offset   *uint32             `json:"offset,omitempty"`
	AuthData *AuthenticationData `json:"authData,omitempty"`
}

type SearchQueryResponse struct {
	TotalMatchesCount uint32                    `json:"totalMatchesCount"`
	Matches           []SearchQueryResponseItem `json:"matches"`
	UserMetaInfo      *SearchQueryUserMetaInfo  `json:"userMetaInfo,omitempty"`
}

type SearchQueryUserMetaInfo struct {
	IsSubscribed          bool   `json:"isSubscribed"`
	LastAccessedTimestamp uint64 `json:"lastAccessedTimestamp"`
}

type SearchQueryResponseItem struct {
	Article ArticleMetaInfo `json:"article"`
}

type SetArticleSubscriptionStatusRequest struct {
	AuthData  AuthenticationData `json:"authData"`
	ArticleId string             `json:"articleId"`
	Subscribe bool               `json:"subscribe"`
}

type GetArticleSubscriptionStatusRequest struct {
	AuthData  AuthenticationData `json:"authData"`
	ArticleId string             `json:"articleId"`
}

type SetSearchQuerySubscriptionStatusRequest struct {
	AuthData  AuthenticationData `json:"authData"`
	Query     string             `json:"query"`
	Subscribe bool               `json:"subscribe"`
}

type GetSearchQuerySubscriptionStatusRequest struct {
	AuthData AuthenticationData `json:"authData"`
	Query    string             `json:"query"`
}

type GetArticlesUpdatesResponse struct {
	UpdatedArticles []ArticleMetaInfo `json:"updatedArticles"`
}

type AccessArticleRequest struct {
	AuthData  AuthenticationData `json:"authData"`
	ArticleId string             `json:"articleId"`
}

type AccessArticleResponse struct {
	Article ArticleMetaInfo `json:"article"`
}

type GetSearchQueriesUpdatesResponse struct {
	SearchQueriesUpdates []GetSearchQueriesUpdatesResponseItem `json:"searchQueriesUpdates"`
}

type GetSearchQueriesUpdatesResponseItem struct {
	Query string `json:"query"`
	// presents if there's at least 1 article that matches the query
	LastNewArticleTimestamp *uint64 `json:"lastNewArticleTimestamp,omitempty"`
	LastAccessedTimestamp   uint64  `json:"lastAccessedTimestamp"`
	NewArticlesCount        *uint32 `json:"newArticlesCount,omitempty"`
}

type GetSearchHistoryResponse struct {
	SearchHistory []GetSearchHistoryResponseItem `json:"searchHistory"`
}

type GetSearchHistoryResponseItem struct {
	Query        string                  `json:"query"`
	UserMetaInfo SearchQueryUserMetaInfo `json:"userMetaInfo"`
}

type GetAccessedArticlesHistoryResponse struct {
	AccessedArticlesHistory []ArticleMetaInfo `json:"accessedArticlesHistory"`
}
