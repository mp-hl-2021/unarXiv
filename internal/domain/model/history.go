package model

type UserSearchAccess struct {
    UserId
    Query     string
    Timestamp uint64
}

type UserArticleAccess struct {
    UserId
    ArticleId
    Timestamp uint64
}

type UserSearchHistory struct {
    UserId
    Queries []string
}
