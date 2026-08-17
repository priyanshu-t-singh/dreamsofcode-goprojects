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

var urlMap = make(map[string]bool)
var deadLinks = make(map[string]bool)
var otherLinks = make(map[string]bool)

func main() {
	url := "http://localhost:8080"
	scrapeLinks(url)

	fmt.Println("-------------------")
	fmt.Println("All Links")
	fmt.Println("-------------------")
	for key := range urlMap {
		fmt.Println(key)
	}

	fmt.Println("-------------------")
	fmt.Println("Other Domain Links")
	fmt.Println("-------------------")
	for key := range otherLinks {
		fmt.Println(key)
	}

	fmt.Println("-------------------")
	fmt.Println("Dead Links")
	fmt.Println("-------------------")
	for key := range deadLinks {
		fmt.Println(key)
	}
}

func get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		deadLinks[url] = true
		return nil, fmt.Errorf("Failed to create resquest: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Go-Fetcher/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		deadLinks[url] = true
		return nil, fmt.Errorf("Failed to fetch the URL: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		deadLinks[url] = true
		return nil, fmt.Errorf("Non-OK HTTP status: %s", resp.Status)
	}

	return resp, nil
}

func markVisited(url string) {
	urlMap[url] = true
}

func scrapeLinks(rawURL string) {
	if _, exists := urlMap[rawURL]; exists {
		fmt.Printf("Url already traversed: %s\n", rawURL)
		return
	}

	if _, exists := deadLinks[rawURL]; exists {
		fmt.Printf("Dead Url already traversed: %s\n", rawURL)
		return
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("failed to parse the url=%q, err=%v\n", rawURL, err)
	}

	markVisited(rawURL)

	resp, err := get(rawURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing the html: %v\n", err)
	}

	traverse(u.Scheme+"://"+u.Host, doc)
}

func traverse(baseUrl string, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key != "href" {
				continue
			}

			fmt.Printf("Anchor tag found -> Href: %s\n", attr.Val)
			href := attr.Val

			if strings.HasPrefix(href, "http") {
				if strings.HasPrefix(href, baseUrl) {
					scrapeLinks(href)
					continue
				}

				if _, exists := urlMap[href]; exists {
					fmt.Printf("Url already traversed: %s\n", href)
					continue
				}

				if _, err := get(href); err != nil {
					fmt.Println(err)
				}
				markVisited(href)
			} else {
				scrapeLinks(strings.TrimSuffix(baseUrl+href, "/"))
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(baseUrl, c)
	}
}
