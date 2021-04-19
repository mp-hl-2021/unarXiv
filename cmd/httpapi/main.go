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
    "database/sql"
    "flag"
    "fmt"
    "net/http"
    "time"

    _ "github.com/lib/pq"
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

    db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", os.Getenv("dbusername"), os.Getenv("dbname")))
    if err != nil {
        panic(err)
    }
    defer db.Close()

    authUsecases := auth.NewUsecases(memory.NewAccountsRepo(db), jwtAuth)
    articleRepo := memory.NewArticleRepo(db)
    historyRepo := memory.NewHistoryRepo()
    artSubsRepo := memory.NewArticleSubscriptionRepo(db)
    searchSubsRepo := memory.NewSearchSubscriptionRepo()
    updatesRepo := repository.NewUpdatesRepoThroughQueries(articleRepo, artSubsRepo, searchSubsRepo, historyRepo)

    if articleDataPath != nil {
        loadArticlesFromFile(articleDataPath, articleRepo)
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

func loadArticlesFromFile(articleDataPath *string, articleRepo *memory.ArticleRepo) {
    var data []byte
    var err error
    if data, err = os.ReadFile(*articleDataPath); err != nil {
        fmt.Printf("Error occurred while reading a file: %v\n", err)
        return
    }
    reader := bufio.NewReader(bytes.NewReader(data))
    _, _ = reader.ReadString('\n') // skip header
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        tokens := strings.Split(strings.Trim(line, "\n"), ";")
        artUrl, err := url.Parse(tokens[0])
        if err != nil {
            fmt.Printf("Failed to parse url: %v\n", err)
            continue
        }

        article := model.Article{
            ArticleMeta: model.ArticleMeta{
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
