package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func extractContent(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var links []string
	z := html.NewTokenizer(res.Body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links, err
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						link := attr.Val
						if strings.HasPrefix(link, "http") {
							links = append(links, link)
						}
						break
					}
				}
			}
		}
	}
}

func crawl(url string, visited map[string]bool, depth int) {
	if depth <= 0 || visited[url] {
		return
	}
	visited[url] = true
	fmt.Printf("Crawling (Depth %d): %s\n", depth, url)

	links, err := extractContent(url)
	if err != nil {
		fmt.Printf("Error extracting content: %v\n", err)
		return
	}

	for _, link := range links {
		crawl(link, visited, depth-1)
	}
}

func main() {
	visited := make(map[string]bool)
	url := "https://go.dev"
	depth := 2
	crawl(url, visited, depth)
}
