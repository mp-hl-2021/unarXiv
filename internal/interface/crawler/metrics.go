package crawler

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	totalURLsVisited = promauto.NewCounter(prometheus.CounterOpts{
		Name: "crawler_total_urls_visited",
		Help: "Number of URLs visited by crawler",
	})
	totalArticlesUpdated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "crawler_total_articles_updated",
		Help: "Number of articles updated by crawler",
	}, []string{"inserted"})
	urlVisitDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "crawler_url_visit_duration_seconds",
		Help: "Duration of a URL visit measured in seconds",
	})
)

func (c *Crawler) RunMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8090", nil); err != nil {
		fmt.Println("Metrics server error:", err)
	}
}
