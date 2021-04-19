package memory

import (
    "github.com/mp-hl-2021/unarXiv/internal/domain"
    "github.com/mp-hl-2021/unarXiv/internal/domain/model"
    "database/sql"
    "fmt"
    "time"

    _ "github.com/lib/pq"
)

type SearchSubscriptionRepo struct {
    db *sql.DB
}

func NewSearchSubscriptionRepo(db *sql.DB) *SearchSubscriptionRepo {
    return &SearchSubscriptionRepo{db: db}
}

func (a *SearchSubscriptionRepo) GetSearchSubscriptions(id model.UserId) ([]string, error) {
    rows, err := a.db.Query(fmt.Sprintf("SELECT Search FROM AccountSearchRelations where UserId = %d;", id))
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    subs := []string{}
    for rows.Next() {
        var query string
        if err := rows.Scan(&query); err != nil {
            panic(err)
        } else {
            subs = append(subs, query)
        }
    }
    return subs, nil
}

func (a *SearchSubscriptionRepo) IsSubscribedForSearch(id model.UserId, query string) (bool, error) {
    rows, err := a.db.Query(fmt.Sprintf("SELECT IsSubscribed FROM AccountSearchRelations WHERE UserId = %s AND Search = '%s';", id, query))
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

func (a *SearchSubscriptionRepo) CreateRelationIfNotExists(userId model.UserId, query string) {
    rows, err := a.db.Query(fmt.Sprintf("SELECT IsSubscribed FROM AccountSearchRelations WHERE UserId = %s AND Search = '%s';", userId, query))
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    relationExists := false
    for rows.Next() {
        relationExists = true
    }
    if !relationExists {
        _, err := a.db.Exec("INSERT INTO AccountSearchRelations (UserId, Search, IsSubscribed, LastAccess) VALUES ($1, $2, false, $3);", userId, query, time.Now())
        if err != nil {
            panic(err)
        }
    }
}


func (a *SearchSubscriptionRepo) SubscribeForSearch(id model.UserId, query string) error {
    if ok, _ := a.IsSubscribedForSearch(id, query); ok {
        return domain.AlreadySubscribed
    }
    a.CreateRelationIfNotExists(id, query)
    _, err := a.db.Exec("UPDATE AccountSearchRelations SET IsSubscribed = true WHERE UserId = $1 AND Search = $2;", id, query)
    if err != nil {
        panic(err)
    }
    return nil
}

func (a *SearchSubscriptionRepo) UnsubscribeFromSearch(id model.UserId, query string) error {
    if ok, _ := a.IsSubscribedForSearch(id, query); ok {
        _, err := a.db.Exec("UPDATE AccountSearchRelations SET IsSubscribed = false WHERE UserId = $1 AND Search = $2;", id, query)
        if err != nil {
            panic(err)
        }
        return nil
    } else {
        return domain.NotSubscribed
    }
}

func (a *SearchSubscriptionRepo) SearchAccessOccurred(id model.UserId, query string) error {
    a.CreateRelationIfNotExists(id, query)
    _, err := a.db.Exec("UPDATE AccountSearchRelations SET LastAccess = $1 WHERE UserId = $2 AND Search = $3;", time.Now(), id, query)
    if err != nil {
        panic(err)
    }
    return nil
}

func (a *SearchSubscriptionRepo) GetSearchLastAccessTimestamp(userId model.UserId, query string) (uint64, error) {
    rows, err := a.db.Query(fmt.Sprintf("SELECT LastAccess FROM AccountSearchRelations WHERE UserId = %s AND Search = '%s';", userId, query))
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

func (a *SearchSubscriptionRepo) GetSearchHistory(userId model.UserId) ([]string, error) {
    rows, err := a.db.Query(fmt.Sprintf("SELECT ArticleId FROM AccountSearchRelations WHERE UserId = %s;", userId))
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    result := []string{}
    for rows.Next() {
        var query string
        if err := rows.Scan(&query); err != nil {
            panic(err)
        } else {
            result = append(result, query)
        }
    }
    return result, nil
}

func (a *SearchSubscriptionRepo) ClearSearchHistory(userId model.UserId) error {
    _, err := a.db.Exec("DELETE FROM AccountSearchRelations WHERE UserId = $1;", userId)
    if err != nil {
        panic(err)
    }
    return nil
}
