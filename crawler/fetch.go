package crawler

import (
	"net/url"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/goware/urlx"

	"github.com/icodeologist/crawly/queue"
	"golang.org/x/net/html"
)

const (
	MAX_DEPTH        = 3
	MAX_LENGTH_LINKS = 100
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
	// the error may be due to wrong type of urls or links
	// ignore them and continue
	if err != nil {
		fmt.Printf("Error GET request : %v\n", err)
	}
	defer resp.Body.Close()
	//parse the response
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		log.Fatalf("Cound not parse the url ", err)
	}
	// Finding all html links and adding it to queue
	list := htmlquery.Find(doc, "//a[@href]")
	for _, n := range list {
		href := htmlquery.SelectAttr(n, "href")
		if href == "" || href[0] == '#' || len(href) >= 10 && href[:10] == "javascript" {
			fmt.Println("Wrong url : ", href)
			continue
		}
		url, err := NormalizeURLs(href)
		if err != nil {
			fmt.Println(url)
			continue
		}
		queue.Que.Enqueue(href)

	}

	// find all h1 to h6 tags
	// and extract its data

	FetchH1TagsData(doc)
	FetchH2TagsData(doc)
	FetchH3TagsData(doc)
	FetchH4TagsData(doc)
	FetchH5TagsData(doc)
	FetchH6TagsData(doc)
}

func FetchH1TagsData(doc *html.Node) {
	// H1 tags
	tags := htmlquery.Find(doc, "//h1")
	for _, n := range tags {
		data := htmlquery.InnerText(n)
		fmt.Println("H1 : ", data)
	}

}

func FetchH2TagsData(doc *html.Node) {
	// H2 tags
	tags := htmlquery.Find(doc, "//h2")
	for _, n := range tags {
		data := htmlquery.InnerText(n)
		fmt.Println("H2 : ", data)
	}
}

func FetchH3TagsData(doc *html.Node) {
	// H3 tags
	tags := htmlquery.Find(doc, "//h3")
	for _, n := range tags {
		data := htmlquery.InnerText(n)
		fmt.Println("H3 : ", data)
	}
}

func FetchH4TagsData(doc *html.Node) {
	// H4 tags
	tags := htmlquery.Find(doc, "//h4")
	for _, n := range tags {
		data := htmlquery.InnerText(n)
		fmt.Println("H4 : ", data)
	}
}

func FetchH5TagsData(doc *html.Node) {
	// H5 tags
	tags := htmlquery.Find(doc, "//h5")
	for _, n := range tags {
		data := htmlquery.InnerText(n)
		fmt.Println("H5 : ", data)
	}
}

func FetchH6TagsData(doc *html.Node) {
	// H6 tags
	tags := htmlquery.Find(doc, "//h6")
	for _, n := range tags {
		data := htmlquery.InnerText(n)
		fmt.Println("H6 : ", data)
	}
}

func FetchPTags(doc *html.Node) {
	ptags := htmlquery.Find(doc, "//p")
	for _, n := range ptags {
		data := htmlquery.InnerText(n)
		fmt.Println("P : ", data)
	}
}

func FetchBoldTags(doc *html.Node) {
	boldTags := htmlquery.Find(doc, "//b")
	for _, n := range boldTags {
		data := htmlquery.InnerText(n)
		fmt.Println("B : ", data)
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
