package main

import (
	"github.com/mp-hl-2021/unarXiv/internal/interface/auth"
	"github.com/mp-hl-2021/unarXiv/internal/interface/httpapi"
	"github.com/mp-hl-2021/unarXiv/internal/interface/repository/memory"
	"github.com/mp-hl-2021/unarXiv/internal/interface/smartUsecases"
	"os"

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
    flag.Parse()

    privateKeyBytes, publicKeyBytes, err := readCryptoKey(*privateKeyPath, *publicKeyPath)
    if err != nil {
        panic(err)
    }

    a, err := auth.NewJwtConfig(privateKeyBytes, publicKeyBytes, 100*time.Minute)
    if err != nil {
        panic(err)
    }

    //unarXivUsecases := dummyUsecases.DummyUsecases{}
    unarXivUsecases := smartUsecases.SmartUsecases{AccountStorage: memory.NewAccountsRepo(), Auth: a}

    httpApi := httpapi.New(&unarXivUsecases)

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
