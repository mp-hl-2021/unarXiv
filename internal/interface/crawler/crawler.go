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
	rows, err := c.db.Query("SELECT COUNT(*) FROM Articles;")
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
	prevArt, err := c.articlesRepo.ArticleById(article.Id)
	if err == domain.ArticleNotFound {
		return true, c.articlesRepo.UpdateArticle(article)
	} else if err != nil {
		return false, err
	}
	differs := prevArt.Title != article.Title ||
		prevArt.Abstract != article.Abstract ||
		len(prevArt.Authors) != len(article.Authors)
	if !differs {
		for i := range prevArt.Authors {
			differs = differs || prevArt.Authors[i] != article.Authors[i]
		}
	}
	if differs {
		return true, c.articlesRepo.UpdateArticle(article)
	}
	return false, nil
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
		urlQueue = urlQueue[1:]
		collector.Visit(u)
		articlesCnt, err = c.getArticlesCount()
	}
	if err != nil {
		return err
	}
	return parseErr
}

func (c *Crawler) responseProcessor(parseErr *error, urlQueue *[]string, cfg *Configuration) func(response *colly.Response) {
	return func(response *colly.Response) {
		originalUrl := response.Request.URL.String()
		dom, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			parseErr = &err
			return
		}
		dom.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			suburl, exists := s.Attr("href")
			if !exists {
				parseErr = &ErrExpectedHref
				return
			}
			if len(*urlQueue) < kURLBacklogSize && strings.Contains(suburl, cfg.RootURL) {
				*urlQueue = append(*urlQueue, suburl)
			}
		})
		if strings.Contains(originalUrl, "abs/") {
			spl := strings.Split(originalUrl, "abs/")
			absId := spl[len(spl)-1]
			if len(absId) < 3 {
				parseErr = &ErrTooShortAbsId
				return
			}
			getElemText := func(class string) string {
				sel := fmt.Sprintf("[class=\"%s\"]", class)
				return strings.Trim(strings.Replace(dom.Find(sel).Text(), "\n", " ", -1), " \t")
			}
			title := getElemText("title mathjax")
			if len(title) == 0 {
				parseErr = &ErrEmptyTitle
				return
			}
			authorsRaw := getElemText("authors")
			if len(authorsRaw) == 0 {
				parseErr = &ErrEmptyAuthors
				return
			}
			authors := strings.Split(authorsRaw, ", ")
			abstract := getElemText("abstract mathjax")
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
			up, err := c.upsertArticle(article)
			if err != nil {
				parseErr = &err
				return
			}
			if up {
				fmt.Println("Upserted article", absId)
			}
		}
	}
}
