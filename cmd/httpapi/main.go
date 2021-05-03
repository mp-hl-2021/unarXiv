package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/mp-hl-2021/unarXiv/internal/interface/auth"
	"github.com/mp-hl-2021/unarXiv/internal/interface/httpapi"
	"github.com/mp-hl-2021/unarXiv/internal/interface/repository/implicitrepos"
	"github.com/mp-hl-2021/unarXiv/internal/interface/repository/postgres"
	"github.com/mp-hl-2021/unarXiv/internal/usecases"
	"net/http"
	"os"
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
	flag.Parse()

	privateKeyBytes, publicKeyBytes, err := readCryptoKey(*privateKeyPath, *publicKeyPath)
	if err != nil {
		panic(err)
	}

	jwtAuth, err := auth.NewJwtConfig(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	dbConnStr := fmt.Sprintf("postgres://%s@db/%s?sslmode=disable", os.Getenv("dbusername"), os.Getenv("dbname"))
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	authUsecases := auth.NewUsecases(postgres.NewAccountsRepo(db), jwtAuth)
	articleRepo := postgres.NewArticleRepo(db)
	artSubsRepo := postgres.NewArticleSubscriptionRepo(db)
	searchSubsRepo := postgres.NewSearchSubscriptionRepo(db)
	updatesRepo := implicitrepos.NewUpdatesRepoThroughQueries(articleRepo, artSubsRepo, searchSubsRepo)

	unarXivUsecases := usecases.NewUsecases(
		authUsecases,
		articleRepo,
		updatesRepo,
		artSubsRepo,
		searchSubsRepo)

	httpApi := httpapi.New(unarXivUsecases)

	httpServer := http.Server{
		Addr:         ":8080",
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
