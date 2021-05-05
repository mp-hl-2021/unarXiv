package postgres

import (
	"database/sql"
	"fmt"
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	// "regexp"
	// "strings"

	_ "github.com/lib/pq"
)

type ArticleRepo struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

func (a *ArticleRepo) ArticleById(id model.ArticleId) (model.Article, error) {
	rows, err := a.db.Query("SELECT Id, Title, Abstract, LastUpdateTimestamp FROM Articles WHERE Id = $1;", id)
	if err != nil {
		return model.Article{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.Id, &article.Title, &article.Abstract, &article.LastUpdateTimestamp); err != nil {
			return model.Article{}, err
		} else {
			rows, err := a.db.Query("SELECT AuthorName FROM AuthorsOfArticles where ArticleId = $1;", id)
			if err != nil {
				return model.Article{}, err
			}
			defer rows.Close()
			for rows.Next() {
				var authorName string
				if err := rows.Scan(&authorName); err != nil {
					return model.Article{}, err
				} else {
					article.Authors = append(article.Authors, authorName)
				}
			}
			return article, nil
		}
	}
	return model.Article{}, domain.ArticleNotFound
}

func (a *ArticleRepo) ArticleMetaById(id model.ArticleId) (model.ArticleMeta, error) {
	article, err := a.ArticleById(id)
	return article.ArticleMeta, err
}

func (a *ArticleRepo) UpdateArticle(article model.Article) error {
	tx, err := a.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	_, err = a.ArticleById(article.Id)
	if err != nil {
		_, err = tx.Exec("INSERT INTO Articles (Id, Title, Abstract, LastUpdateTimestamp, FullDocumentURL) VALUES ($1, $2, $3, $4, $5);", string(article.ArticleMeta.Id), article.ArticleMeta.Title, article.ArticleMeta.Abstract, article.LastUpdateTimestamp, article.FullDocumentURL.String())
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO ArticlesFTS (Id, TextData) VALUES ($1, to_tsvector($2))", string(article.ArticleMeta.Id), fmt.Sprint(article))
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec("UPDATE Articles SET Title = $1, Abstract = $2, LastUpdateTimestamp = $3, FullDocumentURL = $4 WHERE Id = $5;", article.ArticleMeta.Title, article.ArticleMeta.Abstract, article.LastUpdateTimestamp, article.FullDocumentURL.String(), article.ArticleMeta.Id)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE ArticlesFTS SET ArticlesFTS.TextData = to_tsvector($1)", fmt.Sprint(article))
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM AuthorsOfArticles WHERE ArticleId = $1;", article.ArticleMeta.Id)
	if err != nil {
		return err
	}

	for _, authorName := range article.ArticleMeta.Authors {
		_, err = tx.Exec("INSERT INTO AuthorsOfArticles (ArticleId, AuthorName) VALUES ($1, $2);", article.ArticleMeta.Id, authorName)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

const searchQueryTotalMatchesCount = `
SELECT COUNT(*)
FROM ArticlesFTS
WHERE TextData @@ plainto_tsquery($1);
`
const searchQuery = `
SELECT ArticlesFTS.Id, ts_rank(TextData, plainto_tsquery($1))
FROM ArticlesFTS
WHERE TextData @@ plainto_tsquery($1)
ORDER BY ts_rank(TextData, plainto_tsquery($1)) DESC
LIMIT $2 OFFSET $3;
`

func (a *ArticleRepo) Search(query model.SearchQuery, limit uint32) (model.SearchResult, error) { // TODO
	totalMatchesQ := a.db.QueryRow(searchQueryTotalMatchesCount, query.Query)
	var totalMatches int
	if err := totalMatchesQ.Scan(&totalMatches); err != nil {
		return model.SearchResult{}, err
	}
	resp := model.SearchResult{
		TotalMatchesCount: uint32(totalMatches),
		Articles:          nil,
	}
	if limit == 0 {
		limit = 1e9
	}
	rows, err := a.db.Query(searchQuery, query.Query, limit, query.Offset)
	if err != nil {
		return resp, err
	}
	defer rows.Close()
	for rows.Next() {
		var aid string
		var rank float64
		if err := rows.Scan(&aid, &rank); err != nil {
			return resp, err
		}
		article, err := a.ArticleById(model.ArticleId(aid))
		if err != nil {
			return resp, err
		}
		resp.Articles = append(resp.Articles, article.ArticleMeta)
	}
	return resp, nil
}
