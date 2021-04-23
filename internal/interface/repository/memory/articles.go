package memory

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain"
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "database/sql"
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
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        var article model.Article
        if err := rows.Scan(&article.Id, &article.Title, &article.Abstract, &article.LastUpdateTimestamp); err != nil {
            panic(err)
        } else {
            rows, err := a.db.Query("SELECT AuthorName FROM AuthorsOfArticles where ArticleId = $1;", id)
            if err != nil {
                panic(err)
            }
            defer rows.Close()
            for rows.Next() {
                var authorName string
                if err := rows.Scan(&authorName); err != nil {
                    panic(err)
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
    _, err := a.ArticleById(article.Id)
    if err != nil {
        _, err = a.db.Exec("INSERT INTO Articles (Id, Title, Abstract, LastUpdateTimestamp, FullDocumentURL) VALUES ($1, $2, $3, $4, $5);", string(article.ArticleMeta.Id), article.ArticleMeta.Title, article.ArticleMeta.Abstract, article.LastUpdateTimestamp, article.FullDocumentURL.String())
        if err != nil {
            panic(err)
        }
    } else {
        _, err = a.db.Exec("UPDATE Articles SET Title = $1, Abstract = $2, LastUpdateTimestamp = $3, FullDocumentURL = $4 WHERE Id = $5;", article.ArticleMeta.Title, article.ArticleMeta.Abstract, article.LastUpdateTimestamp, article.FullDocumentURL.String(), article.ArticleMeta.Id)
        if err != nil {
            panic(err)
        }
    }

    _, err = a.db.Exec("DELETE FROM AuthorsOfArticles WHERE ArticleId = $1;", article.ArticleMeta.Id)
    if err != nil {
        panic(err)
    }

    for _, authorName := range(article.ArticleMeta.Authors) {
        _, err = a.db.Exec("INSERT INTO AuthorsOfArticles (ArticleId, AuthorName) VALUES ($1, $2);", article.ArticleMeta.Id, authorName)
        if err != nil {
            panic(err)
        }
    }
    return nil
}

func (a *ArticleRepo) Search(query model.SearchQuery, limit uint32) (model.SearchResult, error) { // TODO
    return model.SearchResult{}, nil
    /*
    a.mutex.Lock()
    defer a.mutex.Unlock()
    offset := query.Offset
    result := model.SearchResult{
        TotalMatchesCount: 0,
        Articles:          []model.ArticleMeta{},
    }

    processMatch := func(article model.Article) {
        result.TotalMatchesCount++
        if result.TotalMatchesCount < offset {
            return
        }

        if limit == 0 || result.TotalMatchesCount < limit + offset {
            result.Articles = append(result.Articles, article.ArticleMeta)
        }
    }

    var re *regexp.Regexp
    var err error

    articleContainsMatch := func(a model.Article, r *regexp.Regexp) bool {
        return r.MatchString(a.Title) ||
               r.MatchString(a.Abstract) ||
               r.MatchString(strings.Join(a.Authors, " "))
    }

    if re, err = regexp.Compile(query.Query); err != nil {
        return result, err
    }
    for _, article := range a.articles {
        if articleContainsMatch(article, re) {
            processMatch(article)
        }
    }
    return result, nil
    */
}
