package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "time"

    "github.com/gocolly/colly/v2"
)

func main() {
    fmt.Println("Go go yees")
    c := colly.NewCollector()

    used := make(map[string]bool)
    urlQueue := make([]string, 0)
    articles := make([]string, 0)

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        url := e.Attr("href")
        if !used[url] && strings.Contains(url, "arxiv.org") {
            if strings.Contains(url, "abs") {
                articles = append(articles, url)
            }
            used[url] = true
            urlQueue = append(urlQueue, url)
        }
	})

    c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

    urlQueue = append(urlQueue, "https://arxiv.org/covid19search")

    needArticles := 10

    for len(articles) < needArticles && len(urlQueue) > 0 {
        url := urlQueue[0]
        urlQueue = urlQueue[1:]
        c.Visit(url)
    }

    articles = articles[:needArticles]

    for _, s := range articles {
        fmt.Println("Found article", s)
    }

    f, err := os.Create("data.txt")
    if err != nil {
        panic(err)
    }
    defer f.Close()
    w := bufio.NewWriter(f)

    fmt.Fprintf(w, "Id;Title;Authors;Abstract;LastUpdateTimestamp\n")

    for i, url := range articles {
        fmt.Println("Checking article number", i)
        c2 := colly.NewCollector()
        authors := ""
        c2.OnHTML(".authors", func(e *colly.HTMLElement) {
            authors = e.Text
        })
        title := ""
        c2.OnHTML("[class=\"title mathjax\"]", func(e *colly.HTMLElement) {
            title = e.Text
        })
        abstract := ""
        c2.OnHTML("[class=\"abstract mathjax\"]", func(e *colly.HTMLElement) {
            abstract = strings.Replace(e.Text, "\n", " ", -1)
        })

        c2.Visit(url)
        fmt.Fprintf(w, "%s;%s;%s;%s;%s\n", url, title, authors, abstract, time.Now().String())
    }
}

