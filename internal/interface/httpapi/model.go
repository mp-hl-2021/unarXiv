package httpapi

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "github.com/mp-hl-2021/unarXiv/internal/usecases"
)

type AuthRequest struct {
    Login    string `json:"login"`
    Password string `json:"password"`
}


type ArticleMetaResponse struct {
    Id                  model.ArticleId `json:"article_id"`
    Title               string          `json:"title"`
    Authors             []string        `json:"authors"`
    Abstract            string          `json:"abstract"`
    LastUpdateTimestamp uint64          `json:"last_update"`
}

func renderArticleMeta(article model.ArticleMeta) ArticleMetaResponse {
    return ArticleMetaResponse{
        Id:                  article.Id,
        Title:               article.Title,
        Authors:             article.Authors,
        Abstract:            article.Abstract,
        LastUpdateTimestamp: article.LastUpdateTimestamp,
    }
}

type ArticleResponse struct {
    ArticleMetaResponse `json:"article_meta"`
    FullDocumentURL     string `json:"full_document_url"`
}

func renderArticle(article model.Article) ArticleResponse {
    return ArticleResponse{
        ArticleMetaResponse: renderArticleMeta(article.ArticleMeta),
        FullDocumentURL:     article.FullDocumentURL.String(),
    }
}

type UserArticleSubscriptionResponse struct {
    UserId    model.UserId    `json:"user_id"`
    ArticleId model.ArticleId `json:"article_id"`
}

func renderUserArticleSubscription(subscription model.UserArticleSubscription) UserArticleSubscriptionResponse {
    return UserArticleSubscriptionResponse{
        UserId:    subscription.UserId,
        ArticleId: subscription.ArticleId,
    }
}

type UserSearchSubscriptionResponse struct {
    UserId model.UserId `json:"user_id"`
    Query  string       `json:"query"`
}

func renderUserSearchSubscription(subscription model.UserSearchSubscription) UserSearchSubscriptionResponse {
    return UserSearchSubscriptionResponse{
        UserId: subscription.UserId,
        Query:  subscription.Query,
    }
}

type UserSearchHistoryResponse struct {
    UserId  model.UserId `json:"user_id"`
    Queries []string     `json:"queries"`
}

func renderUserSearchHistory(history model.UserSearchHistory) UserSearchHistoryResponse {
    return UserSearchHistoryResponse{
        UserId:  history.UserId,
        Queries: history.Queries,
    }
}

type UserArticleHistoryResponse struct {
    UserId   model.UserId          `json:"user_id"`
    Articles []ArticleMetaResponse `json:"articles"`
}

func renderUserArticleHistory(history model.UserArticleHistory) UserArticleHistoryResponse {
    articles := make([]ArticleMetaResponse, len(history.Articles))
    for i := range history.Articles {
        articles[i] = renderArticleMeta(history.Articles[i])
    }
    return UserArticleHistoryResponse{
        UserId:   history.UserId,
        Articles: articles,
    }
}

type SearchResultResponse struct {
    TotalMatchesCount uint32                `json:"total_matches"`
    Articles          []ArticleMetaResponse `json:"articles"`
}

func renderSearchResults(result model.SearchResult) SearchResultResponse {
    r := SearchResultResponse{
        TotalMatchesCount: result.TotalMatchesCount,
        Articles:          make([]ArticleMetaResponse, len(result.Articles)),
    }
    for i := range result.Articles {
        r.Articles[i] = renderArticleMeta(result.Articles[i])
    }
    return r
}

type AuthTokenResponse struct {
    Token string `json:"token"`
}

func renderAuthToken(token usecases.AuthToken) AuthTokenResponse {
    return AuthTokenResponse{Token: string(token)}
}
