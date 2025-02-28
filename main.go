package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	fmt.Println("Hello, World!")
	args := os.Args[1:]
	maxConcurrency := 10
	maxPages := 9999
	dl := len(args)
	if dl == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if dl > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	var err error
	// If there are at least 2 arguments, update maxPages
	if dl == 2 {
		maxPages, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("error, not an integer value:", args[1])
			os.Exit(1)
		}
	}
	// If there are 3 arguments, update maxConcurrency
	if dl == 3 {
		maxConcurrency, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("error, not an integer value:", args[1])
			os.Exit(1)
		}
		maxPages, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("error, not an integer value:", args[2])
			os.Exit(1)
		}
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
		maxPages:           maxPages,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency), // Limit do X jednoczesnych watkow (np. 5)
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
