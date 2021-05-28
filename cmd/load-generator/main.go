package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var config = struct {
	address          string
	concurrencyLevel int
}{}

var queries = []string{"plasma", "derivative", "nuclear", "learning", "theorem", "aspect", "entropy"}

func init() {
	address := flag.String("address", "http://localhost:8080", "chat address")
	concurrencyLevel := flag.Int("concurrency", 1, "a number of concurrent requests")
	flag.Parse()

	config.address = *address
	config.concurrencyLevel = *concurrencyLevel
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer func() {
		signal.Stop(ch)
		cancel()
	}()

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	c := client{
		c: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	var wg sync.WaitGroup
	wg.Add(config.concurrencyLevel)
	for i := 0; i < config.concurrencyLevel; i++ {
		go func(i int) {
			err := worker(ctx, c)
			fmt.Printf("worker %d finished, err: %v\n", i, err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("all workers have finished")
}

func worker(ctx context.Context, c client) error {
	for {
		select {
		default:
			var err error
			if rand.Float32() < 0.7 {
				_, err = c.createAccount(ctx, gofakeit.Username()+gofakeit.DigitN(3), gofakeit.Password(true, true, true, false, false, 16))
			} else {
				_, err = c.search(ctx, queries[rand.Intn(len(queries))])
			}
			if err != nil {
				fmt.Println("request failed:", err)
			}
		case <-ctx.Done():
			fmt.Println("leaving worker")
			return ctx.Err()
		}
	}
}

type client struct {
	c http.Client
}

func (c client) createAccount(ctx context.Context, login, password string) (string, error) {
	body := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{login, password}
	s, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/register", bytes.NewReader(s))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create account: %v", resp.Status)
	}
	return resp.Header.Get("Location"), nil
}

func (c client) search(ctx context.Context, query string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.address + "/search/" + query, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to search: %v", resp.Status)
	}
	return resp.Header.Get("Location"), nil
}