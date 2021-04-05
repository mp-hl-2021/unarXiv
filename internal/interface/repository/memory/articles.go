package memory

import (
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"regexp"
	"strings"
	"sync"
)

type ArticleRepo struct {
	articles map[model.ArticleId]model.Article
	mutex    *sync.Mutex
}

func NewArticleRepo() *ArticleRepo {
	return &ArticleRepo{
		articles: make(map[model.ArticleId]model.Article),
		mutex:    &sync.Mutex{},
	}
}

func (a *ArticleRepo) ArticleMetaById(id model.ArticleId) (model.ArticleMeta, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if article, ok := a.articles[id]; !ok {
		return model.ArticleMeta{}, domain.ArticleNotFound
	} else {
		return article.ArticleMeta, nil
	}
}

func (a *ArticleRepo) ArticleById(id model.ArticleId) (model.Article, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if article, ok := a.articles[id]; !ok {
		return model.Article{}, domain.ArticleNotFound
	} else {
		return article, nil
	}
}

func (a *ArticleRepo) UpdateArticle(article model.Article) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.articles[article.Id] = article
	return nil
}

func (a *ArticleRepo) Search(query model.SearchQuery, limit uint32) (model.SearchResult, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	offsetLeft := query.Offset
	result := model.SearchResult{
		TotalMatchesCount: 0,
		Articles:          []model.ArticleMeta{},
	}
	processMatch := func(article model.Article) {
		result.TotalMatchesCount++
		if offsetLeft > 0 {
			offsetLeft--
		} else if limit == 0 || uint32(len(result.Articles)) < limit {
			result.Articles = append(result.Articles, article.ArticleMeta)
		}
	}
	var re *regexp.Regexp
	var err error
	if re, err = regexp.Compile(query.Query); err != nil {
		return result, err
	}
	for _, article := range a.articles {
		if re.MatchString(article.Title) ||
			re.MatchString(article.Abstract) ||
			re.MatchString(strings.Join(article.Authors, " ")) {
			processMatch(article)
		}
	}
	return result, nil
}
