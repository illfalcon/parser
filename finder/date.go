package finder

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
)

var months = []string{
	"январь", "января", "февраль", "февраля", "март", "марта", "апрель", "апреля", "май", "мая",
	"июнь", "июня", "июль", "июля", "август", "августа", "сентябрь", "сентября", "октябрь", "октября",
	"ноябрь", "ноября", "декабрь", "декабря",
}

func containsMonth(text string) bool {
	for _, month := range months {
		if strings.Contains(text, month) {
			return true
		}
	}
	return false
}

func ContainsDate() (bool, error) {
	resp, err := http.Get("http://kdobru.ru/info/events/events_309.html")
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return false, fmt.Errorf("getting %s: %s", "http://kdobru.ru/info/events/events_309.html", resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("script").Remove()
	root := doc.Get(0)
	var nodesWithData []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && containsMonth(n.Data) {
			nodesWithData = append(nodesWithData, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(root)
	sel := &goquery.Selection{}
	sel = sel.AddNodes(nodesWithData...)
	newSel := &goquery.Selection{}
	sel.Each(func(i int, selection *goquery.Selection) {
		newSel = newSel.AddSelection(selection.ParentsUntil("div").Parent().Contents())
	})
	fmt.Println(newSel.Text())
	return true, nil
}
