package usecases

import (
    "github.com/mp-hl-2021/unarXiv/internal/usecases/subscriptions"
)

type Interface interface {
    AuthInterface
    ArticleInterface
    SearchInterface
    HistoryInterface
    subscriptions.ArticleInterface
    subscriptions.SearchInterface
}
