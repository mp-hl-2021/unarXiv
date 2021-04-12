package main

import (
	"bufio"
	"bytes"
	"github.com/mp-hl-2021/unarXiv/internal/domain/model"
	"github.com/mp-hl-2021/unarXiv/internal/interface/auth"
	"github.com/mp-hl-2021/unarXiv/internal/interface/httpapi"
	"github.com/mp-hl-2021/unarXiv/internal/interface/repository"
	"github.com/mp-hl-2021/unarXiv/internal/interface/repository/memory"
	"github.com/mp-hl-2021/unarXiv/internal/interface/utils"
	"github.com/mp-hl-2021/unarXiv/internal/usecases"
	"net/url"
	"os"
	"strings"

	"flag"
	"fmt"
	"net/http"
	"time"
)

func readCryptoKey(privateKeyPath string, publicKeyPath string) (privateKeyBytes []byte, publicKeyBytes []byte, err error) {
    privateKeyBytes, err = os.ReadFile(privateKeyPath)
    if err != nil {
        return
    }
    publicKeyBytes, err = os.ReadFile(publicKeyPath)
    return
}

func main() {
    privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
    publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
    articleDataPath := flag.String("articleDataPath", "data.txt", "file path")
    flag.Parse()

    privateKeyBytes, publicKeyBytes, err := readCryptoKey(*privateKeyPath, *publicKeyPath)
    if err != nil {
        panic(err)
    }

    jwtAuth, err := auth.NewJwtConfig(privateKeyBytes, publicKeyBytes, 100*time.Minute)
    if err != nil {
        panic(err)
    }

    authUsecases := auth.NewUsecases(memory.NewAccountsRepo(), jwtAuth)
    articleRepo := memory.NewArticleRepo()
    historyRepo := memory.NewHistoryRepo()
    artSubsRepo := memory.NewArticleSubscriptionRepo()
    searchSubsRepo := memory.NewSearchSubscriptionRepo()
    updatesRepo := repository.NewUpdatesRepoThroughQueries(articleRepo, artSubsRepo, searchSubsRepo, historyRepo)

    if articleDataPath != nil {
		if data, err := os.ReadFile(*articleDataPath); err != nil {
			fmt.Printf("Error occurred while reading a file: %v\n", err)
 		} else {
			reader := bufio.NewReader(bytes.NewReader(data))
			_, _ = reader.ReadString('\n')  // skip header
			for {
				if line, err := reader.ReadString('\n'); err != nil {
					break
				} else {
					tokens := strings.Split(strings.Trim(line, "\n"), ";")
					artUrl, err := url.Parse(tokens[0])
					if err != nil {
						fmt.Printf("Failed to parse url: %v\n", err)
						continue
					}

					article := model.Article{
						ArticleMeta:     model.ArticleMeta{
							Id:                  model.ArticleId(tokens[0]),
							Title:               tokens[1],
							Authors:             strings.Split(tokens[2], ", "),
							Abstract:            tokens[3],
							LastUpdateTimestamp: utils.Uint64Time(time.Now()),
						},
						FullDocumentURL: *artUrl,
					}
					if err = articleRepo.UpdateArticle(article); err != nil {
						fmt.Printf("Failed to put an article into the repo: %v\n", err)
					} else {
						fmt.Printf("Added article %v to the repo!\n", article.Id)
					}
				}
			}
		}
	} else {
		fmt.Println("article data path is not specified -> zero articles loaded")
	}

    unarXivUsecases := usecases.NewUsecases(
		authUsecases,
    	articleRepo,
    	historyRepo,
    	updatesRepo,
    	artSubsRepo,
    	searchSubsRepo)

    httpApi := httpapi.New(unarXivUsecases)

    httpServer := http.Server{
        Addr:         "localhost:8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,

        Handler: httpApi.Router(),
    }
    fmt.Println("Listening on :8080")
    err = httpServer.ListenAndServe()
    if err != nil {
        panic(err)
    }
}
