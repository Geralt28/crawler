package main

import (
	"fmt"
	"io"
	"net/http"
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
