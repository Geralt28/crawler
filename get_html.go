package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	url, err := normalizeURL(rawURL)
	if err != nil {
		fmt.Println("error: bad url format")
		os.Exit(1)
	}
	//fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error: could not get page")
		return "", err
	}
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
	defer res.Body.Close()
	return string(body), nil
}

func crawlPage(rawBaseURL string, pages map[string]int) {
	normBaseURL, err := normalizeURL(rawBaseURL)
	if err != nil {
		fmt.Println("could not normalise base url:", rawBaseURL)
		return
	}
	baseURL, err := url.Parse(normBaseURL)
	if err != nil {
		fmt.Println("could not parse base url:", normBaseURL)
		return
	}
	var currentURL string
	listURL := []string{normBaseURL}
	//pages[normBaseURL] = 1
	for {
		if len(listURL) == 0 {
			break
		}
		currentURL = listURL[0]
		listURL = listURL[1:]
		//normCurrentURL, err := normalizeURL(rawCurrentURL)
		//if err != nil {
		//	fmt.Println("could not normalise current url:", rawBaseURL)
		//	continue
		//}
		cURL, err := url.Parse(currentURL)
		if err != nil {
			fmt.Println("could not parse current url:", rawBaseURL)
			continue
		}
		if cURL.Host != baseURL.Host {
			//fmt.Println("different hosts:", currentURL)
			continue
		}
		if _, ok := pages[currentURL]; ok {
			pages[currentURL] += 1
			continue
		} else {
			pages[currentURL] = 1
		}
		html, err := getHTML(currentURL)
		if err != nil {
			continue
		}
		fmt.Println("...crawling now:", currentURL)
		newListURL, err := getURLsFromHTML(html, rawBaseURL)
		if err != nil {
			fmt.Println("error: could not get urls from html")
			continue
		}
		listURL = append(listURL, newListURL...)
	}

}
