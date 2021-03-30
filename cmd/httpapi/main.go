package main

import (
    "fmt"
    "github.com/mp-hl-2021/unarXiv/internal/interface/smartUsecases"
    "github.com/mp-hl-2021/unarXiv/internal/interface/httpapi"
    "net/http"
    "time"
)

func main() {
    //unarXivUsecases := dummyUsecases.DummyUsecases{}
    unarXivUsecases := smartUsecases.SmartUsecases{}

    httpApi := httpapi.New(&unarXivUsecases)

    httpServer := http.Server{
        Addr:         "localhost:8080",
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,

        Handler: httpApi.Router(),
    }
    fmt.Println("Listening on :8080")
    err := httpServer.ListenAndServe()
    if err != nil {
        panic(err)
    }
}
