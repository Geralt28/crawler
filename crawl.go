package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	maxPages           int
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Stop adding new pages if the limit is reached (for stopping delayed adding)
	if len(cfg.pages) >= cfg.maxPages {
		return false
	}

	if _, exists := cfg.pages[normalizedURL]; exists {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		fmt.Println("error: could not get page")
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode > 400 {
		return "", fmt.Errorf("error status code: %d", res.StatusCode)
	}
	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("error: invalid text/html format: %s", res.Header.Get("Content-Type"))
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error: could not read body:", err)
		return "", err
	}
	return string(body), nil
}

func (cfg *config) crawlPage(normCurrentURL string) {

	defer cfg.wg.Done()
	cfg.concurrencyControl <- struct{}{}        // Acquire a slot
	defer func() { <-cfg.concurrencyControl }() // Release the slot when done

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	cURL, err := url.Parse(normCurrentURL)
	if err != nil {
		fmt.Println("could not parse URL:", normCurrentURL)
		return
	}
	if cURL.Host != cfg.baseURL.Host {
		return // Pomin zewnetrzne linki
	}
	if !cfg.addPageVisit(normCurrentURL) {
		return // Pomin juz odwiedzone strony
	}
	fmt.Println("...crawling now:", normCurrentURL)
	html, err := getHTML(normCurrentURL)
	if err != nil {
		return
	}
	newListURL, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		fmt.Println("error: could not get URLs from HTML")
		return
	}
	for _, newURL := range newListURL {
		cfg.wg.Add(1)
		go cfg.crawlPage(newURL) // Spawn goroutine
	}
}
