package main

import (
	"log"

	"github.com/illfalcon/parser/internal/collyparser"
)

const url string = "http://kdobru.ru/moodle/"

func main() {
	err := collyparser.CrawlAndParse(url, 2)
	if err != nil {
		log.Fatal(err)
	}
	//db.Prepare()
	//frontend.Start()
}
