package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

type PageResult struct {
	URL        string
	StatusCode int
	NewLinks   []string
	IsExternal bool
	Err        error
}

func main() {
	startURLStr := "http://localhost:8080"
	baseURL, err := url.Parse(startURLStr)
	if err != nil {
		log.Fatal(err)
	}

	resultsChan := make(chan PageResult)
	visited := make(map[string]int)

	// Track active workers to know when to exit
	activeWorkers := 0

	// Start the first worker
	activeWorkers++
	visited[startURLStr] = 0 // 0 means "currently processing"
	go processPage(startURLStr, baseURL, false, resultsChan)

	// Orchestrator Loop
	for activeWorkers > 0 {
		res := <-resultsChan
		activeWorkers--

		if res.Err != nil {
			visited[res.URL] = -1
		} else {
			visited[res.URL] = res.StatusCode
		}

		for _, link := range res.NewLinks {
			if _, alreadyVisited := visited[link]; !alreadyVisited {
				visited[link] = 0
				activeWorkers++

				isExternal := !strings.HasPrefix(link, baseURL.Scheme+"://"+baseURL.Host)
				go processPage(link, baseURL, isExternal, resultsChan)
			}
		}
	}

	fmt.Println("-------------------")
	fmt.Println("Working Links")
	fmt.Println("-------------------")
	for url, statusCode := range visited {
		if statusCode > 0 && statusCode < 400 {
			fmt.Printf("%d -> %s\n", statusCode, url)
		}
	}

	fmt.Println("-------------------")
	fmt.Println("Dead Links")
	fmt.Println("-------------------")
	for url, statusCode := range visited {
		if statusCode >= 400 || statusCode == -1 {
			fmt.Printf("%d -> %s\n", statusCode, url)
		}
	}
}

func processPage(targetURL string, baseURL *url.URL, isExternal bool, out chan<- PageResult) {
	res := PageResult{URL: targetURL, IsExternal: isExternal}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		res.Err = fmt.Errorf("failed to create request: %v", err)
		out <- res
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-Fetcher/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		res.Err = fmt.Errorf("failed to fetch: %v", err)
		out <- res
		return
	}
	defer resp.Body.Close()

	res.StatusCode = resp.StatusCode

	// Skip HTML parsing for external links or dead links
	if isExternal || resp.StatusCode != http.StatusOK {
		out <- res
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		res.Err = fmt.Errorf("error parsing html: %v", err)
		out <- res
		return
	}

	res.NewLinks = extractLinks(baseURL, doc)
	out <- res
}

func extractLinks(baseURL *url.URL, n *html.Node) []string {
	var links []string

	var traverse func(n *html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					parsedHref, err := url.Parse(attr.Val)
					if err == nil {
						resolved := baseURL.ResolveReference(parsedHref).String()
						resolved = strings.TrimSuffix(strings.Split(resolved, "#")[0], "/")
						links = append(links, resolved)
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return links
}
