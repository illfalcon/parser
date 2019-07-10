package finder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//type dateFinder struct {
//}
//
//func (m dateFinder) Match(node *html.Node) bool {
//	return node.Type == html.TextNode && containsMonth(node.Data)
//}
//
//func (m dateFinder) MatchAll(*html.Node) []*html.Node {
//
//}
//
//func (m dateFinder) Filter([]*html.Node) []*html.Node {
//
//}

var months = []string{
	" январь", " января", " февраль", " февраля", " март", " марта", " апрель", " апреля", " май", " мая",
	" июнь", " июня", " июль", " июля", " август", " августа", " сентябрь", " сентября", " октябрь", " октября",
	" ноябрь", " ноября", " декабрь", " декабря",
}

func containsMonth(text string) bool {
	for _, month := range months {
		if strings.Contains(text, month) {
			return true
		}
	}
	return false
}

func isInlineTag(tag string) bool {
	return tag == "em" || tag == "strong" || tag == "span"
}

func WriteDivsWithDate(resp *http.Response, num int) (bool, error) {
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("unable to read file")
	}
	htmlText := string(htmlBytes)
	htmlText = strings.ReplaceAll(htmlText, "<em>", "")
	htmlText = strings.ReplaceAll(htmlText, "<strong>", "")
	htmlText = strings.ReplaceAll(htmlText, "<span>", "")
	htmlText = strings.ReplaceAll(htmlText, "</em>", "")
	htmlText = strings.ReplaceAll(htmlText, "</strong>", "")
	htmlText = strings.ReplaceAll(htmlText, "</span>", "")
	reg := regexp.MustCompile(`<!--[^>]*-->`)
	htmlText = reg.ReplaceAllString(htmlText, "")
	r := strings.NewReader(htmlText)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("script").Remove()

	//root := doc.Get(0)
	//var nodesWithData []*html.Node
	//var f func(*html.Node)
	//f = func(n *html.Node) {
	//	if n.Type == html.TextNode && containsMonth(n.Data) {
	//		nodesWithData = append(nodesWithData, n)
	//		log.Println(n.Data)
	//	}
	//	for c := n.FirstChild; c != nil; c = c.NextSibling {
	//		f(c)
	//	}
	//}
	//f(root)
	f := func(i int, selection *goquery.Selection) bool {
		return containsMonth(selection.Text())
	}
	g := func(i int, selection *goquery.Selection) bool {
		var b bool
		selection.Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) == "#text" {
				if f(i, s) {
					b = true
				}
			}
		})
		return b
	}
	sel := doc.Find("*").FilterFunction(g)
	//	NotFunction(func(i int, selection *goquery.Selection) bool {
	//	return f(i, selection.Children())
	//})
	//for _, n := range newSel.Nodes {
	//	if n.Type == html.TextNode && containsMonth(n.Data) {
	//		nodesWithData = append(nodesWithData, n)
	//		log.Println(n.Data)
	//	}
	//}
	//sel := doc.Selection.FindNodes(nodesWithData...)
	file, err := os.Create("/home/elgreco/Desktop/contents1/" + strconv.Itoa(num) + ".txt")
	if err != nil {
		return false, fmt.Errorf("unable to create file %s: %s", resp.Request.URL.Path+".txt", err)
	}
	defer file.Close()
	//html.Render(file, doc.Get(0))
	sel.Each(func(i int, selection *goquery.Selection) {
		_, err = file.WriteString(strings.TrimSpace(selection.Closest("div").Contents().Text()) +
			"\n_______________________\n")
		if err != nil {
			log.Printf("unable to write to file %s: %s", resp.Request.URL.Path+".txt", err)
		}
	})
	//if newSel.Length() != 0 {
	//	file, err := os.Create("/home/elgreco/Desktop/contents/" + strconv.Itoa(num) + ".txt")
	//	if err != nil {
	//		return false, fmt.Errorf("unable to create file %s: %s", resp.Request.URL.Path+".txt", err)
	//	}
	//	defer file.Close()
	//	newSel.Each(func(i int, selection *goquery.Selection) {
	//		_, err = file.WriteString(strings.TrimSpace(selection.Text()))
	//		if err != nil {
	//			log.Printf("unable to write to file %s: %s", resp.Request.URL.Path+".txt", err)
	//		}
	//	})
	//}
	return true, nil
}
