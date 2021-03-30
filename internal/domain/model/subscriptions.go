package model

type UserArticleSubscription struct {
    UserId
    ArticleId
}

type UserSearchSubscription struct {
    UserId
    Query string
}
