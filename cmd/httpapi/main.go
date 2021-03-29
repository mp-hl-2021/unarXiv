package main

import (
	"github.com/mp-hl-2021/unarXiv/internal/api"
	"github.com/mp-hl-2021/unarXiv/internal/usecases"
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
