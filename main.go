package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	fmt.Println("Hello, World!")
	args := os.Args[1:]
	dl := len(args)
	if dl == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if dl > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	rawBaseUrl := args[0]
	normBaseURL, err := normalizeURL(rawBaseUrl)
	if err != nil {
		fmt.Println("could not normalise base URL:", rawBaseUrl)
		os.Exit(1)
	}
	baseURL, err := url.Parse(normBaseURL)
	if err != nil {
		fmt.Println("error: invalid base URL format:", normBaseURL)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 10), // Limit do X jednoczesnych watkow (np. 5)
		wg:                 &sync.WaitGroup{},
	}

	fmt.Println("starting crawl of:", normBaseURL)
	cfg.wg.Add(1)
	go cfg.crawlPage(normBaseURL)
	cfg.wg.Wait()

	fmt.Println("\n****** List of pages ******")
	for url, count := range cfg.pages {
		fmt.Println(url, ":", count)
	}
}
