package main

import (
    "github.com/mp-hl-2021/unarXiv/internal/interface/smartUsecases"
    "github.com/mp-hl-2021/unarXiv/internal/interface/httpapi"
    "github.com/mp-hl-2021/unarXiv/internal/interface/accountstorage"
    "github.com/mp-hl-2021/unarXiv/internal/interface/auth"

    "fmt"
    "net/http"
    "time"
    "flag"
    "io/ioutil"
)

func main() {
    privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)

	a, err := auth.NewJwt(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

    //unarXivUsecases := dummyUsecases.DummyUsecases{}
    unarXivUsecases := smartUsecases.SmartUsecases{AccountStorage: accountstorage.NewMemory(), Auth: a}

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
