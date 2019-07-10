package main

import (
	"github.com/illfalcon/parser/internal/crawler"
	"github.com/illfalcon/parser/internal/parser"
)

func main() {
	//db.Prepare()
	crawler.Crawl()
	parser.Parse()
}
