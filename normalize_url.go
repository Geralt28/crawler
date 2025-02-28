package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(inputURL string) (string, error) {
	if len(inputURL) == 0 {
		return "", fmt.Errorf("pusty adress")
	}
	// Remove spaces additional "/" and make it lower
	inputURL = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(inputURL), "/"))
	// Parse the input URL, if there is no "http://"" try to add it
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	// Remove the scheme (http/https)
	normalized := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path
	// Remove trailing slashes
	normalized = strings.TrimSuffix(normalized, "/")
	return normalized, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	n, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		fmt.Println("err: could not read html nodes")
		return nil, err
	}
	var urls []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url := strings.TrimSpace(attr.Val)
					if len(url) > 0 {
						if url[0] == '/' {
							url = strings.TrimRight(rawBaseURL, "/") + url
						}
					}
					normalizedURL, err := normalizeURL(url)
					if err == nil {
						urls = append(urls, normalizedURL)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return urls, nil
}
