package main

import (
	"database/sql"
	"fmt"
	"github.com/mp-hl-2021/unarXiv/internal/interface/crawler"
	"github.com/mp-hl-2021/unarXiv/internal/interface/repository/postgres"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	dbConnStr := fmt.Sprintf("postgres://%s@db/%s?sslmode=disable", os.Getenv("dbusername"), os.Getenv("dbname"))
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	articleRepo := postgres.NewArticleRepo(db)

	c := crawler.NewCrawler(db, articleRepo)

	for {
		cfg, err := c.GetConfiguration()
		if err != nil {
			panic(err)
		}
		if err := c.CrawlArticles(cfg); err != nil {
			panic(err)
		}
		time.Sleep(time.Minute)
	}
}
