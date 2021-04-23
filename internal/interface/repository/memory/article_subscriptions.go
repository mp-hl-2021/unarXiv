package memory

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain"
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "database/sql"
    "time"

    _ "github.com/lib/pq"
)

type ArticleSubscriptionRepo struct {
    db *sql.DB
}

func NewArticleSubscriptionRepo(db *sql.DB) *ArticleSubscriptionRepo {
    return &ArticleSubscriptionRepo{db: db}
}

func (a *ArticleSubscriptionRepo) GetArticleSubscriptions(id model.UserId) ([]model.ArticleId, error) {
    rows, err := a.db.Query("SELECT ArticleId FROM AccountArticleRelations where UserId = $1;", id)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    subs := []model.ArticleId{}
    for rows.Next() {
        var articleId model.ArticleId
        if err := rows.Scan(&articleId); err != nil {
            panic(err)
        } else {
            subs = append(subs, articleId)
        }
    }
    return subs, nil
}

func (a *ArticleSubscriptionRepo) IsSubscribedForArticle(userId model.UserId, articleId model.ArticleId) (bool, error) {
    rows, err := a.db.Query("SELECT IsSubscribed FROM AccountArticleRelations WHERE UserId = $1 AND ArticleId = $2;", userId, articleId)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        var isSubscribed bool
        if err := rows.Scan(&isSubscribed); err != nil {
            panic(err)
        } else {
            return isSubscribed, nil
        }
    }
    return false, nil
}

func (a *ArticleSubscriptionRepo) CreateRelationIfNotExists(userId model.UserId, articleId model.ArticleId) {
    rows, err := a.db.Query("SELECT IsSubscribed FROM AccountArticleRelations WHERE UserId = $1 AND ArticleId = $2;", userId, articleId)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    relationExists := false
    for rows.Next() {
        relationExists = true
    }
    if !relationExists {
        _, err := a.db.Exec("INSERT INTO AccountArticleRelations (UserId, ArticleId, IsSubscribed, LastAccess) VALUES ($1, $2, false, $3);", userId, articleId, time.Now())
        if err != nil {
            panic(err)
        }
    }
}

func (a *ArticleSubscriptionRepo) SubscribeForArticle(id model.UserId, articleId model.ArticleId) error {
    if ok, _ := a.IsSubscribedForArticle(id, articleId); ok {
        return domain.AlreadySubscribed
    }
    a.CreateRelationIfNotExists(id, articleId)
    _, err := a.db.Exec("UPDATE AccountArticleRelations SET IsSubscribed = true WHERE UserId = $1 AND ArticleID = $2;", id, articleId)
    if err != nil {
        panic(err)
    }
    return nil
}

func (a *ArticleSubscriptionRepo) UnsubscribeFromArticle(id model.UserId, articleId model.ArticleId) error {
    if ok, _ := a.IsSubscribedForArticle(id, articleId); ok {
        _, err := a.db.Exec("UPDATE AccountArticleRelations SET IsSubscribed = false WHERE UserId = $1 AND ArticleID = $2;", id, articleId)
        if err != nil {
            panic(err)
        }
        return nil
    } else {
        return domain.NotSubscribed
    }
}

func (a *ArticleSubscriptionRepo) ArticleAccessOccurred(id model.UserId, articleId model.ArticleId) error {
    a.CreateRelationIfNotExists(id, articleId)
    _, err := a.db.Exec("UPDATE AccountArticleRelations SET LastAccess = $1 WHERE UserId = $2 AND ArticleID = $3;", time.Now(), id, articleId)
    if err != nil {
        panic(err)
    }
    return nil
}

func (a *ArticleSubscriptionRepo) GetArticleLastAccessTimestamp(userId model.UserId, articleId model.ArticleId) (uint64, error) {
    rows, err := a.db.Query("SELECT LastAccess FROM AccountArticleRelations WHERE UserId = $1 AND ArticleId = $2;", userId, articleId)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        var lastAccess uint64
        if err := rows.Scan(&lastAccess); err != nil {
            panic(err)
        } else {
            return lastAccess, nil
        }
    }
    return 0, domain.NeverAccessed
}

func (a *ArticleSubscriptionRepo) GetArticleHistory(userId model.UserId) ([]model.ArticleId, error) {
    rows, err := a.db.Query("SELECT ArticleId FROM AccountArticleRelations WHERE UserId = $1;", userId)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    result := []model.ArticleId{}
    for rows.Next() {
        var articleId model.ArticleId
        if err := rows.Scan(&articleId); err != nil {
            panic(err)
        } else {
            result = append(result, articleId)
        }
    }
    return result, nil
}

func (a *ArticleSubscriptionRepo) ClearArticleHistory(userId model.UserId) error {
    _, err := a.db.Exec("DELETE FROM AccountArticleRelations WHERE UserId = $1;", userId)
    if err != nil {
        panic(err)
    }
    return nil
}
