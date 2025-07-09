package main

import (
	"github.com/icodeologist/crawly/crawler"
	"github.com/icodeologist/crawly/queue"
)

func main() {
	url := "https://wikipedia.com"
	queue.Que.Enqueue(url)
	crawler.QueueLogic()

}
