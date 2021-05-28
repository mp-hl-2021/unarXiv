package crawler

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/signal"
	"net/http"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"github.com/mp-hl-2021/unarXiv/internal/domain"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/domain/repository"
	"github.com/mp-hl-2021/unarXiv/internal/interface/utils"
	"strings"
	"time"
)

var (
	ErrNoConfigs     = fmt.Errorf("no configurations found")
	ErrExpectedHref  = fmt.Errorf("expected attribute href but it wasn't found")
	ErrTooShortAbsId = fmt.Errorf("too short absId")
	ErrEmptyTitle    = fmt.Errorf("parsed empty title")
	ErrEmptyAuthors  = fmt.Errorf("parsed empty authors")
	ErrEmptyQueue    = fmt.Errorf("no more urls to crawl")

	chBuff                 = 10
	downloadURLConcurrency = 2
	parseHTMLConcurrency   = 2
	putArticleLConcurrency = 1
	putURLConcurrency      = 1
)


type Crawler struct {
	db		   *sql.DB
	articlesRepo repository.ArticleRepo
}

func NewCrawler(db *sql.DB, articlesRepo repository.ArticleRepo) *Crawler {
	return &Crawler{db: db, articlesRepo: articlesRepo}
}

type Configuration struct {
	RootURL			 string
	DesiredArticleCount int
}

func (c *Crawler) getUnvisitedURLs() ([]string, error) {
	rows, err := c.db.Query("SELECT URL FROM CrawlStatus WHERE Visited = false LIMIT 100;") // todo: fix magic number. probably in the next life
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	urls := make([]string, 0, 100) // todo: fix magic number. probably in the next life
	for rows.Next() {
		url := ""
		err := rows.Scan(&url)
		if err != nil {
			return []string{}, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (c *Crawler) addURLToQueue(url string) error {
	_, err := c.db.Exec("INSERT INTO CrawlStatus (URL, Visited) VALUES ($1, false) ON CONFLICT DO NOTHING;", url)
	return err
}

func (c *Crawler) dbVisitURL(url string) error {
	_, err := c.db.Exec("UPDATE CrawlStatus SET Visited = true where URL = $1;", url)
	return err
}

func (c *Crawler) dbUpdateURLInfo(url string, HTTPStatus int) error {
	_, err := c.db.Exec("UPDATE CrawlStatus SET LastAccess = $1, LastHTTPStatus = $2 where URL = $3;", utils.Uint64Time(time.Now()), HTTPStatus, url)
	return err
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

func (c *Crawler) getURLFromDB(ctx context.Context, URLChan chan<- string) error {
	var urls []string
	var err error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if len(urls) == 0 {
				urls, err = c.getUnvisitedURLs()
				if err != nil {
					return err
				}
				if len(urls) == 0 {
					time.Sleep(time.Second)
					continue
				}
			}
			url := urls[0]
			urls = urls[1:]
			err = c.dbVisitURL(url)
			if err != nil {
				return err
			}
			URLChan <- url
		}
	}
}

func (c *Crawler) downloadURL(ctx context.Context, URLChan <-chan string, HTMLChan chan<- *http.Response) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case url := <-URLChan:
			fmt.Println("Visiting", url)
			response, err := http.Get(url)
			if err != nil {
				return err
			}
			totalURLsVisited.Inc()
			err = c.dbUpdateURLInfo(url, response.StatusCode)
			if err != nil {
				return err
			}
			HTMLChan <- response
			time.Sleep(time.Second)
		}
	}
}

func (c *Crawler) parseHTML(ctx context.Context, cfg *Configuration, HTMLChan <-chan *http.Response, ArticleChan chan<- model.Article, NewURLChan chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case response := <-HTMLChan:
			body, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			dom, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
			if err != nil {
				return err
			}
			err = c.collectUrls(dom, cfg, NewURLChan)
			if err != nil {
				return err
			}
			if strings.Contains(response.Request.URL.String(), "/abs/") {
				article, err := c.parseArticle(response, dom)
				if err != nil {
					return err
				}
				ArticleChan <- article
			}
			response.Body.Close()
		}
	}
}

func (c *Crawler) putURLToDB(ctx context.Context, NewURLChan <-chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case url := <-NewURLChan:
			err := c.addURLToQueue(url)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Crawler) putArticleToDB(ctx context.Context, ArticleChan <-chan model.Article) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case article := <-ArticleChan:
			up, err := c.upsertArticle(article)
			if err != nil {
				return err
			}
			if up {
				fmt.Println("Upserted article", article.Id)
			}
		}
	}
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

	ctx, cancel := context.WithCancel(context.Background())
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	defer func() {
		signal.Stop(osChan)
		cancel()
	}()

	URLChan := make(chan string, chBuff)
	HTMLChan := make(chan *http.Response, chBuff)
	ArticleChan := make(chan model.Article, chBuff)
	NewURLChan := make(chan string, chBuff)

	go func() {
		for {
			select {
			case <-osChan:
				cancel()
			case <-ctx.Done():
				return
			}
		}
	}()

	var gwg sync.WaitGroup
	gwg.Add(1)
	go func(out chan<- string) {
		err := c.getURLFromDB(ctx, out)
		fmt.Fprintf(os.Stderr, "URLGetter 1 stopped, reason: %s\n", err)
		gwg.Done()
		cancel()
	}(URLChan)

	var dwg sync.WaitGroup
	dwg.Add(downloadURLConcurrency)
	for i := 0; i < downloadURLConcurrency; i++ {
		go func(in <-chan string, out chan<- *http.Response) {
			err := c.downloadURL(ctx, in, out)
			fmt.Fprintf(os.Stderr, "URLDownloader %d stopped, reason: %s\n", i, err)
			dwg.Done()
			cancel()
		}(URLChan, HTMLChan)
	}

	var parseWG sync.WaitGroup
	parseWG.Add(parseHTMLConcurrency)
	for i := 0; i < parseHTMLConcurrency; i++ {
		go func(in <-chan *http.Response, outArticle chan<- model.Article, outURL chan<- string) {
			err := c.parseHTML(ctx, &cfg, in, outArticle, outURL)
			fmt.Fprintf(os.Stderr, "HTMLParser %d stopped, reason: %s\n", i, err)
			parseWG.Done()
			cancel()
		}(HTMLChan, ArticleChan, NewURLChan)
	}

	var putUrlWG sync.WaitGroup
	putUrlWG.Add(putURLConcurrency)
	for i := 0; i < putURLConcurrency; i++ {
		go func(in <-chan string) {
			err := c.putURLToDB(ctx, in)
			fmt.Fprintf(os.Stderr, "URLPutter %d stopped, reason: %s\n", i, err)
			putUrlWG.Done()
			cancel()
		}(NewURLChan)
	}

	var putArticleWG sync.WaitGroup
	putArticleWG.Add(putArticleLConcurrency)
	for i := 0; i < putArticleLConcurrency; i++ {
		go func(in <-chan model.Article) {
			err := c.putArticleToDB(ctx, in)
			fmt.Fprintf(os.Stderr, "ArticlePutter %d stopped, reason: %s\n", i, err)
			putArticleWG.Done()
			cancel()
		}(ArticleChan)
	}

	gwg.Wait()
	dwg.Wait()
	parseWG.Wait()
	putUrlWG.Wait()
	putArticleWG.Wait()

	return ctx.Err()
}

func (c *Crawler) parseArticle(response *http.Response, dom *goquery.Document) (model.Article, error) {
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
			Id:				  model.ArticleId(absId),
			Title:			   title,
			Authors:			 authors,
			Abstract:			abstract,
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

func (c *Crawler) collectUrls(dom *goquery.Document, cfg *Configuration, NewURLChan chan<- string) error {
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
		if strings.HasPrefix(suburl, "/") {
			suburl = cfg.RootURL + suburl[1:]
		}
		if strings.Contains(suburl, cfg.RootURL) && !strings.Contains(suburl, "/pdf/") && !strings.Contains(suburl, "/ps/") {
			NewURLChan <- suburl
		}
	})
	return err
}
