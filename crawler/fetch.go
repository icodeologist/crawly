package crawler

import (
	"net/url"

	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/goware/urlx"
	"log"
	"net/http"
	"time"

	"github.com/icodeologist/crawly/queue"
)

const (
	MAX_DEPTH        = 3
	MAX_LENGTH_LINKS = 20
)

// map the links to bool to keep track of
// duplicates and further crawling

func HandleCrashErrors(err error) {
	log.Fatalf("Program crashed : %s\n", err)
}

func HandleNonCrashableErrors(err error) {
	log.Printf("Unexpected error : %e\n", err)
}

//having client to have control over redirect and timeouts

type Client struct {
	Client *http.Client
}

// reusable constructor
func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// map to keep check of visited links
// this prevent from revisiting the duplicated links
var visitedLinks = make(map[string]bool)

// given link crawl the link and put the extracted data to database
// return links and add them to queue and crawl them back
func Crawl(url string, depth int) {
	if depth >= 3 {
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error : Could not connect to the server %v\n", err)
	}
	defer resp.Body.Close()
	//parse the response
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Fatalf("Cound not parse the url ", err)
	}
	list := htmlquery.Find(doc, "//a[@href]")
	for _, n := range list {
		href := htmlquery.SelectAttr(n, "href")
		prefix := "https:"
		queue.Que.Enqueue(prefix + href)
	}
}

// this function handles visiting the extracted links from the crawl function
// these links will be checked for duplicates
// then we pass these links to crawl further more links
// With keeping our constraints in check
//
// CONSTRAINTS
// max_depth <= 3
// max_pages <= 20 for now
func QueueLogic() {
	depth := 0
	for !queue.Que.IsEmpty() {
		link, ok := queue.Que.Dequeue()
		if !ok {
			return
		}
		if link == "" {
			continue
		}
		if visitedLinks[link] == true {
			continue
		}
		if len(visitedLinks) >= MAX_LENGTH_LINKS {
			fmt.Printf("Maximum Link quota reached %v Want : %v Have", MAX_LENGTH_LINKS, len(visitedLinks))
			for l, _ := range visitedLinks {
				fmt.Println("l", l)

			}
			fmt.Println("Total visited sites : ", len(visitedLinks))
			return
		}
		visitedLinks[link] = true
		fmt.Printf("Crawling %v\n", link)
		Crawl(link, depth+1)
	}

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
