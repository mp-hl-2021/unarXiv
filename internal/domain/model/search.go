package model

type SearchQuery struct {
    Query  string
    Offset uint32
}

type SearchResult struct {
    TotalMatchesCount uint32
    Articles          []ArticleMeta
}
