package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/gocolly/colly/v2"
    "github.com/tkanos/gonfig"
)

type CrawlerConfiguration struct {
    CrawlDomain      string
    CrawlStart       string
    LocalStoragePath string
    NeedArticles     int
}

func GetCrawlerConfiguration(cfgPath string) CrawlerConfiguration {
    cfg := CrawlerConfiguration{}
    err := gonfig.GetConf(cfgPath, &cfg)
    if err != nil {
        panic(err)
    }
    return cfg
}

func CrawlArticles(cfg CrawlerConfiguration) []string {
    fmt.Println("Start crawling")
    c := colly.NewCollector()

    used := make(map[string]bool)
    urlQueue := make(chan string, 1000)
    articles := make([]string, 0)

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        url := e.Attr("href")
        if !used[url] && strings.Contains(url, cfg.CrawlDomain) {
            if strings.Contains(url, "abs") {
                articles = append(articles, url)
            }
            used[url] = true
            urlQueue <- url
        }
    })
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })


    urlQueue <- cfg.CrawlStart
    for len(articles) < cfg.NeedArticles && len(urlQueue) > 0 {
        c.Visit(<-urlQueue)
    }

    return articles[:cfg.NeedArticles]
}

func GetFieldFromHTML(c *colly.Collector, saveTo *string, className string) {
    c.OnHTML(fmt.Sprintf("[class=\"%s\"]", className), func (e *colly.HTMLElement) {
        *saveTo = strings.Replace(e.Text, "\n", " ", -1)
    })
}

func SaveArticlesInfo(articles []string, cfg CrawlerConfiguration) {
    f, err := os.Create(cfg.LocalStoragePath)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    defer w.Flush()

    fmt.Fprintf(w, "Id;Title;Authors;Abstract;LastUpdateTimestamp\n")
    for i, url := range articles {
        fmt.Println("Saving info about article number", i)
        c2 := colly.NewCollector()
        authors := ""
        GetFieldFromHTML(c2, &authors, "authors")
        title := ""
        GetFieldFromHTML(c2, &title, "title mathjax")
        abstract := ""
        GetFieldFromHTML(c2, &abstract, "abstract mathjax")

        c2.Visit(url)
        fmt.Fprintf(w, "%s;%s;%s;%s;%s\n", url, title, authors, abstract, time.Now().String())
    }
}

func main() {
    cfg := GetCrawlerConfiguration("cmd/crawler/crawlercfg.json")

    articles := CrawlArticles(cfg)

    for _, s := range articles {
        fmt.Println("Found article", s)
    }

    SaveArticlesInfo(articles, cfg)
}
