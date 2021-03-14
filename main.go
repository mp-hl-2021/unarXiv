package main

import (
	"github.com/mp-hl-2021/unarXiv/core"
	"github.com/mp-hl-2021/unarXiv/server"
	"net/http"
	"time"
)

func main() {
	unarXivAPI := core.DummyUnarXivAPI{}

	service := server.NewServer(unarXivAPI)

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
