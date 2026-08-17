package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

var (
	resultMap = make(map[string]int)
	mapMutex  sync.Mutex
)

func main() {
	url := "http://localhost:8080"
	var wg sync.WaitGroup

	wg.Add(1)
	go scrapeLinks(url, &wg)
	wg.Wait()

	fmt.Println("-------------------")
	fmt.Println("Working Links")
	fmt.Println("-------------------")
	for url, statusCode := range resultMap {
		if statusCode < 400 {
			fmt.Printf("%d -> %s\n", statusCode, url)
		}
	}

	fmt.Println("-------------------")
	fmt.Println("Dead Links")
	fmt.Println("-------------------")
	for url, statusCode := range resultMap {
		if statusCode >= 400 {
			fmt.Printf("%d -> %s\n", statusCode, url)
		}
	}
}

func markIfNew(u string) bool {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if _, exists := resultMap[u]; exists {
		return false
	}

	resultMap[u] = 0
	return true
}

func updateStatus(u string, statusCode int) {
	mapMutex.Lock()
	defer mapMutex.Unlock()
	resultMap[u] = statusCode
}

func get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create resquest: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-Fetcher/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the URL: %v", err)
	}

	updateStatus(url, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-OK HTTP status: %s", resp.Status)
	}

	return resp, nil
}

func scrapeLinks(rawURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	if !markIfNew(rawURL) {
		return
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("failed to parse the url=%q, err=%v\n", rawURL, err)
		return
	}

	resp, err := get(rawURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing the html: %v\n", err)
		return
	}

	traverse(u, doc, wg)
}

func traverse(baseURL *url.URL, n *html.Node, wg *sync.WaitGroup) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}

			fmt.Printf("Anchor tag found -> Href: %s\n", attr.Val)
			parsedHref, err := url.Parse(attr.Val)
			if err != nil {
				continue
			}
			resolvedURL := baseURL.ResolveReference(parsedHref).String()
			resolvedURL = strings.TrimSuffix(strings.Split(resolvedURL, "#")[0], "/")

			if strings.HasPrefix(resolvedURL, baseURL.Scheme+"://"+baseURL.Host) {
				wg.Add(1)
				go scrapeLinks(resolvedURL, wg)
			} else {
				if markIfNew(resolvedURL) {
					wg.Add(1)
					go func(extURL string) {
						defer wg.Done()
						if _, err := get(extURL); err != nil {
							fmt.Println(err)
						}
					}(resolvedURL)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(baseURL, c, wg)
	}
}
