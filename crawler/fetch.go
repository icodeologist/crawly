package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error : Parsing html file by go query generated this error %v\n", err)
	}
	// TODO: extract the data here
	// extractData(link)
	// TODO: Put the data to databse
	// make the struct and put the data to database

	// Gather all the links from the downloaded html body
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		// TODO: Check domain
		if exists && strings.HasPrefix(link, "https") && !visitedLinks[link] {
			queue.Que.Enqueue(link)
		}
	})
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
			log.Fatalf("Queue is empty")
			return
		}
		if link == "" {
			continue
		}
		if visitedLinks[link] == true {
			continue
		}
		if len(visitedLinks) >= MAX_LENGTH_LINKS {
			log.Fatalf("Maximum Link quota reached %v Want : %v Have", MAX_LENGTH_LINKS, len(visitedLinks))
			return
		}
		visitedLinks[link] = true
		fmt.Printf("Crawling %v\n", link)
		Crawl(link, depth+1)
	}
}
