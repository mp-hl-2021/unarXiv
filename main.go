package main

import (
"github.com/mp-hl-2021/unarXiv/api"
"github.com/mp-hl-2021/unarXiv/server"
"net/http"
"time"
)

func main() {
	useCases := api.DummyUseCases{}

	service := server.NewApi(useCases)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
