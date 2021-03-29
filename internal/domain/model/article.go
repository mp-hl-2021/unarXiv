package model

import "net/url"

type ArticleId string
type ArticleMeta struct {
    Id                  ArticleId
    Title               string
    Authors             []string
    Abstract            string
    LastUpdateTimestamp uint64
}

type Article struct {
    ArticleMeta

    FullDocumentURL url.URL
}
