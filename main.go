package main

import (
	"github.com/mp-hl-2021/unarXiv/core"
	"github.com/mp-hl-2021/unarXiv/api"
	"net/http"
	"time"
)

func main() {
	unarXivUseCases := core.DummyUnarXivAPI{}

	unarXivApi := api.NewApi(unarXivUseCases)

	httpServer := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: unarXivApi.Router(),
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
