package repository

import (
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/domain/repository"
)

type UpdatesRepoThroughQueries struct {
	articleRepo     repository.ArticleRepo
	articleSubsRepo repository.ArticleSubscriptionRepo
	searchSubsRepo  repository.SearchSubscriptionRepo
}

func NewUpdatesRepoThroughQueries(
	articleRepo repository.ArticleRepo,
	articleSubsRepo repository.ArticleSubscriptionRepo,
	searchSubsRepo repository.SearchSubscriptionRepo) *UpdatesRepoThroughQueries {
	return &UpdatesRepoThroughQueries{
		articleRepo:     articleRepo,
		articleSubsRepo: articleSubsRepo,
		searchSubsRepo:  searchSubsRepo,
	}
}

func (u *UpdatesRepoThroughQueries) GetArticleSubscriptionsUpdates(id model.UserId) ([]model.ArticleMeta, error) {
	subs, err := u.articleSubsRepo.GetArticleSubscriptions(id)
	if err != nil {
		return nil, err
	}
	var result []model.ArticleMeta
	for _, articleId := range subs {
		articleMeta, err := u.articleRepo.ArticleMetaById(articleId)
		if err == domain.ArticleNotFound {
			continue
		} else if err != nil {
			return nil, err
		}
		lastAccessTimestamp, err := u.articleSubsRepo.GetArticleLastAccessTimestamp(id, articleId)
		if err != nil && err != domain.NeverAccessed {
			return nil, err
		} else if lastAccessTimestamp < articleMeta.LastUpdateTimestamp || err == domain.NeverAccessed {
			result = append(result, articleMeta)
		}
	}
	return result, nil
}

/*
Командир демонстрирует солдатам новый танк.
- Вот, товарищи бойцы, это новый секретный танк. Петров.
- Я!
- Подними танк.
Петров тужится, пыжится, не поднять.
- Не поднять.
- Сидоров, помоги Петрову.
Пытаются вдвоем, та же ситуация.
- Не поднять.
- Иванов, помогай.
Пыхтят втроем. Поднять не могут.
- Никак не поднять, товарищ командир!
- А чего вы ожидали? Сорок шесть тонн!
*/
func (u *UpdatesRepoThroughQueries) GetSearchSubscriptionsUpdates(id model.UserId) ([]string, error) {
	subs, err := u.searchSubsRepo.GetSearchSubscriptions(id)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, query := range subs {
		articles, err := u.articleRepo.Search(model.SearchQuery{
			Query:  query,
			Offset: 0,
		}, 0)
		if err != nil {
			return nil, err
		}
		lastArticleUpdateTimestamp := uint64(0)
		for _, articleMeta := range articles.Articles {
			if lastArticleUpdateTimestamp < articleMeta.LastUpdateTimestamp {
				lastArticleUpdateTimestamp = articleMeta.LastUpdateTimestamp
			}
		}
		lastAccessTimestamp, err := u.searchSubsRepo.GetSearchLastAccessTimestamp(id, query)
		if err != nil && err != domain.NeverAccessed {
			return nil, err
		} else if lastAccessTimestamp < lastArticleUpdateTimestamp || err == domain.NeverAccessed {
			result = append(result, query)
		}
	}
	return result, nil
}
