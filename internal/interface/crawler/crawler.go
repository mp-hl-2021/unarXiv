package crawler

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/domain/repository"
	"github.com/mp-hl-2021/unarXiv/internal/interface/utils"
	"strings"
	"time"
)

const kURLBacklogSize = 1000

var (
	ErrNoConfigs     = fmt.Errorf("no configurations found")
	ErrExpectedHref  = fmt.Errorf("expected attribute href but it wasn't found")
	ErrTooShortAbsId = fmt.Errorf("too short absId")
	ErrEmptyTitle    = fmt.Errorf("parsed empty title")
	ErrEmptyAuthors  = fmt.Errorf("parsed empty authors")
)

type Crawler struct {
	db           *sql.DB
	articlesRepo repository.ArticleRepo
}

func NewCrawler(db *sql.DB, articlesRepo repository.ArticleRepo) *Crawler {
	return &Crawler{db: db, articlesRepo: articlesRepo}
}

type Configuration struct {
	RootURL             string
	DesiredArticleCount int
}

func (c *Crawler) GetConfiguration() (Configuration, error) {
	rows, err := c.db.Query("SELECT RootURL, DesiredArticleCount FROM CrawlerConfig;")
	if err != nil {
		return Configuration{}, err
	}
	defer rows.Close()
	for rows.Next() {
		cfg := Configuration{}
		err := rows.Scan(&cfg.RootURL, &cfg.DesiredArticleCount)
		return cfg, err
	}
	return Configuration{}, ErrNoConfigs
}

func (c *Crawler) getArticlesCount() (int, error) {
	rows, err := c.db.Query("SELECT n_live_tup FROM pg_stat_all_tables WHERE relname = 'articles';")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var cnt int
		err = rows.Scan(&cnt)
		return cnt, err
	}
	return 0, fmt.Errorf("unexpected query result")
}

func (c *Crawler) upsertArticle(article model.Article) (bool, error) {
	prevArticleState, err := c.articlesRepo.ArticleById(article.Id)
	if err == domain.ArticleNotFound {
		err = c.articlesRepo.UpdateArticle(article)
		if err != nil {
			return false, err
		}
		totalArticlesUpdated.WithLabelValues("1").Inc()
		return true, nil
	}
	if err != nil {
		return false, err
	}
	if article.Equals(prevArticleState) {
		return false, nil
	}
	err = c.articlesRepo.UpdateArticle(article)
	if err != nil {
		return false, err
	}
	totalArticlesUpdated.WithLabelValues("0").Inc()
	return true, err
}

func (c *Crawler) CrawlArticles(cfg Configuration) error {
	fmt.Println("Crawling...")
	collector := colly.NewCollector()
	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       time.Second / 4,
		RandomDelay: time.Second / 4,
	})
	if err != nil {
		return err
	}

	urlQueue := []string{cfg.RootURL}
	var parseErr error

	collector.OnResponse(c.responseProcessor(&parseErr, &urlQueue, &cfg))
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	articlesCnt, err := c.getArticlesCount()
	for err == nil && articlesCnt < cfg.DesiredArticleCount && len(urlQueue) > 0 && parseErr == nil {
		u := urlQueue[0]
		totalURLsVisited.Inc()
		urlQueue = urlQueue[1:]
		timeStart := time.Now()
		collector.Visit(u)
		urlVisitDuration.Observe(time.Now().Sub(timeStart).Seconds())
		articlesCnt, err = c.getArticlesCount()
	}
	if err != nil {
		return err
	}
	return parseErr
}

func (c *Crawler) responseProcessor(parseErr *error, urlQueue *[]string, cfg *Configuration) func(response *colly.Response) {
	return func(response *colly.Response) {
		dom, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			parseErr = &err
			return
		}
		err = c.collectUrls(dom, urlQueue, cfg)
		if err != nil {
			parseErr = &err
			return
		}
		if strings.Contains(response.Request.URL.String(), "abs/") {
			article, err := c.parseArticle(response, dom)
			if err != nil {
				parseErr = &err
				return
			}
			up, err := c.upsertArticle(article)
			if err != nil {
				parseErr = &err
				return
			}
			if up {
				fmt.Println("Upserted article", article.Id)
			}
		}
	}
}

func (c *Crawler) parseArticle(response *colly.Response, dom *goquery.Document) (model.Article, error) {
	originalUrl := response.Request.URL.String()
	absId, err := c.extractArticleId(originalUrl)
	if err != nil {
		return model.Article{}, err
	}
	title := c.getElemTextByClass(dom, "title mathjax")
	if len(title) == 0 {
		return model.Article{}, ErrEmptyTitle
	}
	authorsRaw := c.getElemTextByClass(dom, "authors")
	if len(authorsRaw) == 0 {
		return model.Article{}, ErrEmptyAuthors
	}
	authors := strings.Split(authorsRaw, ", ")
	abstract := c.getElemTextByClass(dom, "abstract mathjax")
	article := model.Article{
		ArticleMeta: model.ArticleMeta{
			Id:                  model.ArticleId(absId),
			Title:               title,
			Authors:             authors,
			Abstract:            abstract,
			LastUpdateTimestamp: utils.Uint64Time(time.Now()),
		},
		FullDocumentURL: *response.Request.URL,
	}
	return article, nil
}

func (c *Crawler) getElemTextByClass(dom *goquery.Document, class string) string {
	sel := fmt.Sprintf("[class=\"%s\"]", class)
	return strings.Trim(strings.Replace(dom.Find(sel).Text(), "\n", " ", -1), " \t")
}

func (c *Crawler) extractArticleId(originalUrl string) (string, error) {
	spl := strings.Split(originalUrl, "abs/")
	absId := spl[len(spl)-1]
	if len(absId) < 3 {
		return "", ErrTooShortAbsId
	}
	return absId, nil
}

func (c *Crawler) collectUrls(dom *goquery.Document, urlQueue *[]string, cfg *Configuration) error {
	var err error
	dom.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if err != nil {
			return
		}
		suburl, exists := s.Attr("href")
		if !exists {
			err = ErrExpectedHref
			return
		}
		if len(*urlQueue) < kURLBacklogSize && strings.Contains(suburl, cfg.RootURL) {
			*urlQueue = append(*urlQueue, suburl)
		}
	})
	return err
}
