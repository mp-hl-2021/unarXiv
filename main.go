package main

import (
	"github.com/mp-hl-2021/unarXiv/core"
	"github.com/mp-hl-2021/unarXiv/server"
	"net/http"
	"time"
)

func main() {
	unarXivAPI := core.DummyUnarXivAPI{}

	unarXivServer := server.NewServer(unarXivAPI)

	httpServer := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: unarXivServer.Router(),
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
