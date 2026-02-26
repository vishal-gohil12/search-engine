package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
	"golang.org/x/net/html"
)

func normalizeURL(baseURL, href string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	normalized := base.ResolveReference(ref)
	normalized.Fragment = ""
	return normalized.String(), nil
}

func extractContent(url string) ([]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "GoCrawler/1.0")

	res, err := client.Get(url)
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
			if z.Err() == io.EOF {
				return links, nil
			}
			return links, z.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						link := attr.Val
						links = append(links, link)
						break
					}
				}
			}
		}
	}
}

func crawl(baseURL string, visited map[string]bool, depth int, group *robotstxt.Group) {
	if depth <= 0 || visited[baseURL] {
		return
	}

	if !group.Test(baseURL) {
		fmt.Println("Blocked by robots.txt:", baseURL)
		return
	}

	visited[baseURL] = true
	fmt.Printf("Crawling (Depth %d): %s\n", depth, baseURL)

	time.Sleep(500 * time.Millisecond)

	links, err := extractContent(baseURL)
	if err != nil {
		fmt.Printf("Error extracting content: %v\n", err)
		return
	}

	for _, link := range links {
		normalizedLink, err := normalizeURL(baseURL, link)

		if err != nil {
			fmt.Printf("Error normalizing URL: %v\n", err)
			continue
		}

		parsedBase, _ := url.Parse(baseURL)
		parsedNew, _ := url.Parse(normalizedLink)

		if parsedBase.Host == parsedNew.Host {
			crawl(normalizedLink, visited, depth-1, group)
		}

	}
}

func main() {
	visited := make(map[string]bool)
	startURL := "https://go.dev"
	depth := 2

	res, err := http.Get(startURL + "/robots.txt")
	if err != nil {
		fmt.Printf("Error fetching robots.txt: %v\n", err)
		return
	}
	defer res.Body.Close()

	robotsdata, err := robotstxt.FromResponse(res)
	if err != nil {
		fmt.Printf("Error parsing robots.txt: %v\n", err)
		return
	}

	group := robotsdata.FindGroup("*")

	crawl(startURL, visited, depth, group)
}
