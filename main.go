package main

import (
	"github.com/mp-hl-2021/unarXiv/api"
	"github.com/mp-hl-2021/unarXiv/usecases"
	"net/http"
	"time"
)

func main() {
	unarXivUseCases := usecases.DummyUnarXivAPI{}

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
