package main

import (
	"fmt"
	"net/http"
	"sync"
)

type HttpResult struct {
	Url        string
	StatusCode int
}

// 4
func makeRequest(ch chan<- HttpResult, url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("url was error: %v", err)
	}
	ch <- HttpResult{Url: url, StatusCode: resp.StatusCode}
}

// 5
func collect(ch <-chan HttpResult) {
	for msg := range ch {
		fmt.Printf("%s -> %d\n", msg.Url, msg.StatusCode)
	}
}

func main() {
        // 1
	urlChan := make(chan HttpResult,1)
	var wg sync.WaitGroup

	urls := []string{"https://rogerwelin.github.io/",
		"https://golang.org/",
		"https://news.ycombinator.com/",
		"https://www.google.se/shouldbe404",
		"https://www.cpan.org/"}
        // 2
	go collect(urlChan)

        // 3
	for _, url := range urls {
		wg.Add(1)
		go makeRequest(urlChan, url, &wg)
	}

       // 6
	wg.Wait()
	close(urlChan)
}
