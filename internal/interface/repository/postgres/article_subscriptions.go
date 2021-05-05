package postgres

import (
	"database/sql"
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/interface/utils"
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
		return nil, err
	}
	defer rows.Close()
	subs := []model.ArticleId{}
	for rows.Next() {
		var articleId model.ArticleId
		if err := rows.Scan(&articleId); err != nil {
			return nil, err
		} else {
			subs = append(subs, articleId)
		}
	}
	return subs, nil
}

func (a *ArticleSubscriptionRepo) IsSubscribedForArticle(userId model.UserId, articleId model.ArticleId) (bool, error) {
	rows, err := a.db.Query("SELECT IsSubscribed FROM AccountArticleRelations WHERE UserId = $1 AND ArticleId = $2;", userId, articleId)
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

func (a *ArticleSubscriptionRepo) createRelationIfNotExists(userId model.UserId, articleId model.ArticleId) error {
	rows, err := a.db.Query("SELECT IsSubscribed FROM AccountArticleRelations WHERE UserId = $1 AND ArticleId = $2;", userId, articleId)
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
			"INSERT INTO AccountArticleRelations (UserId, ArticleId, IsSubscribed, LastAccess) VALUES ($1, $2, false, $3);",
			userId, articleId, utils.Uint64Time(time.Now()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ArticleSubscriptionRepo) SubscribeForArticle(id model.UserId, articleId model.ArticleId) error {
	ok, err := a.IsSubscribedForArticle(id, articleId)
	if err != nil {
		return err
	}
	if ok {
		return domain.AlreadySubscribed
	}
	err = a.createRelationIfNotExists(id, articleId)
	if err != nil {
		return err
	}
	_, err = a.db.Exec("UPDATE AccountArticleRelations SET IsSubscribed = true WHERE UserId = $1 AND ArticleID = $2;", id, articleId)
	return err
}

func (a *ArticleSubscriptionRepo) UnsubscribeFromArticle(id model.UserId, articleId model.ArticleId) error {
	ok, err := a.IsSubscribedForArticle(id, articleId)
	if err != nil {
		return err
	}
	if ok {
		_, err := a.db.Exec("UPDATE AccountArticleRelations SET IsSubscribed = false WHERE UserId = $1 AND ArticleID = $2;", id, articleId)
		return err
	} else {
		return domain.NotSubscribed
	}
}

func (a *ArticleSubscriptionRepo) ArticleAccessOccurred(id model.UserId, articleId model.ArticleId) error {
	err := a.createRelationIfNotExists(id, articleId)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(
		"UPDATE AccountArticleRelations SET LastAccess = $1 WHERE UserId = $2 AND ArticleID = $3;",
		utils.Uint64Time(time.Now()), id, articleId)
	return err
}

func (a *ArticleSubscriptionRepo) GetArticleLastAccessTimestamp(userId model.UserId, articleId model.ArticleId) (uint64, error) {
	rows, err := a.db.Query("SELECT LastAccess FROM AccountArticleRelations WHERE UserId = $1 AND ArticleId = $2;", userId, articleId)
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

func (a *ArticleSubscriptionRepo) GetArticleHistory(userId model.UserId) ([]model.ArticleId, error) {
	rows, err := a.db.Query("SELECT ArticleId FROM AccountArticleRelations WHERE UserId = $1;", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.ArticleId{}
	for rows.Next() {
		var articleId model.ArticleId
		if err := rows.Scan(&articleId); err != nil {
			return nil, err
		} else {
			result = append(result, articleId)
		}
	}
	return result, nil
}

func (a *ArticleSubscriptionRepo) ClearArticleHistory(userId model.UserId) error {
	_, err := a.db.Exec("DELETE FROM AccountArticleRelations WHERE UserId = $1;", userId)
	if err != nil {
		return err
	}
	return nil
}
