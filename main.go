package main

import (
	"fmt"

	"net/url"

	"github.com/goware/urlx"
	"github.com/icodeologist/crawly/crawler"
	"github.com/icodeologist/crawly/queue"
)

func main() {
	url := "https://example.com"
	queue.Que.Enqueue(url)
	fmt.Println("Crawler has started")
	crawler.QueueLogic()

}

// Normalize the urls
func NormalizeURLs(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("Error : Parsing url %v\n", err)
	}
	fmt.Println("U", u)
	normalizedString, err := urlx.Normalize(u)
	if err != nil {
		return "", err
	} else {
		return normalizedString, nil
	}
}
