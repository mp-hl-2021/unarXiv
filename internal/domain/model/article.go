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

func (a Article) Equals(b Article) bool {
	if a.Id != b.Id ||
		a.Title != b.Title ||
		a.Abstract != b.Abstract ||
		a.FullDocumentURL.String() != b.FullDocumentURL.String() ||
		len(a.Authors) != len(b.Authors) {
		return false
	}
	for i := range a.Authors {
		if a.Authors[i] != b.Authors[i] {
			return false
		}
	}
	return true
}
