package postgres

import (
	"database/sql"
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/interface/utils"
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
	rows, err := a.db.Query("SELECT Search FROM AccountSearchRelations where UserId = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subs := []string{}
	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			return nil, err
		} else {
			subs = append(subs, query)
		}
	}
	return subs, nil
}

func (a *SearchSubscriptionRepo) IsSubscribedForSearch(id model.UserId, query string) (bool, error) {
	rows, err := a.db.Query("SELECT IsSubscribed FROM AccountSearchRelations WHERE UserId = $1 AND Search = $2;", id, query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var isSubscribed bool
		err := rows.Scan(&isSubscribed)
		return isSubscribed, err
	}
	return false, nil
}

func (a *SearchSubscriptionRepo) createRelationIfNotExists(userId model.UserId, query string) error {
	rows, err := a.db.Query("SELECT IsSubscribed FROM AccountSearchRelations WHERE UserId = $1 AND Search = $2;", userId, query)
	if err != nil {
		return err
	}
	defer rows.Close()
	relationExists := false
	for rows.Next() {
		relationExists = true
	}
	if !relationExists {
		_, err := a.db.Exec(
			"INSERT INTO AccountSearchRelations (UserId, Search, IsSubscribed, LastAccess) VALUES ($1, $2, false, $3);",
			userId, query, utils.Uint64Time(time.Now()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *SearchSubscriptionRepo) SubscribeForSearch(id model.UserId, query string) error {
	if ok, _ := a.IsSubscribedForSearch(id, query); ok {
		return domain.AlreadySubscribed
	}
	err := a.createRelationIfNotExists(id, query)
	if err != nil {
		return err
	}
	_, err = a.db.Exec("UPDATE AccountSearchRelations SET IsSubscribed = true WHERE UserId = $1 AND Search = $2;", id, query)
	if err != nil {
		return err
	}
	return nil
}

func (a *SearchSubscriptionRepo) UnsubscribeFromSearch(id model.UserId, query string) error {
	if ok, _ := a.IsSubscribedForSearch(id, query); ok {
		_, err := a.db.Exec("UPDATE AccountSearchRelations SET IsSubscribed = false WHERE UserId = $1 AND Search = $2;", id, query)
		return err
	} else {
		return domain.NotSubscribed
	}
}

func (a *SearchSubscriptionRepo) SearchAccessOccurred(id model.UserId, query string) error {
	err := a.createRelationIfNotExists(id, query)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(
		"UPDATE AccountSearchRelations SET LastAccess = $1 WHERE UserId = $2 AND Search = $3;",
		utils.Uint64Time(time.Now()), id, query)
	return err
}

func (a *SearchSubscriptionRepo) GetSearchLastAccessTimestamp(userId model.UserId, query string) (uint64, error) {
	rows, err := a.db.Query("SELECT LastAccess FROM AccountSearchRelations WHERE UserId = $1 AND Search = $2;", userId, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var lastAccess uint64
		err := rows.Scan(&lastAccess)
		return lastAccess, err
	}
	return 0, domain.NeverAccessed
}

func (a *SearchSubscriptionRepo) GetSearchHistory(userId model.UserId) ([]string, error) {
	rows, err := a.db.Query("SELECT ArticleId FROM AccountSearchRelations WHERE UserId = $1;", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			return nil, err
		} else {
			result = append(result, query)
		}
	}
	return result, nil
}

func (a *SearchSubscriptionRepo) ClearSearchHistory(userId model.UserId) error {
	_, err := a.db.Exec("DELETE FROM AccountSearchRelations WHERE UserId = $1;", userId)
	return err
}
