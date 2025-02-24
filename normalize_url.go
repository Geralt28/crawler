package main

import (
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	// Remove spaces additional "/" and make it lower
	inputURL = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(inputURL), "/"))
	// Parse the input URL, if there is no "http://"" try to add it
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		if parsedURL.Scheme == "" {
			parsedURL, err = url.Parse("http://" + inputURL)
			if err != nil {
				return "", err
			}
		}
	}
	// Remove the scheme (http/https)
	normalized := parsedURL.Host + parsedURL.Path
	// Remove trailing slashes
	normalized = strings.TrimSuffix(normalized, "/")
	return normalized, nil
}
