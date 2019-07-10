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

	"github.com/illfalcon/parser/watson"

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

var invitations = []string{
	"Приглашаем", "приглашаем", "Состоится", "состоится", "пройдет", "Пройдет", "пройдёт", "Пройдёт", "открывается",
	"Открывается", "откроется", "откроется", "проводим", "Проводим", "Проведем", "проведем", "Проведём", "проведём",
	"приглашаются", "Приглашаются", "проведет", "проведёт", "Проведет", "Проведёт",
}

func containsMonth(text string) bool {
	for _, month := range months {
		if strings.Contains(text, month) {
			return true
		}
	}
	return false
}

func containsInvitation(text string) bool {
	for _, inv := range invitations {
		if strings.Contains(text, inv) {
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
	f := func(i int, selection *goquery.Selection) bool {
		return containsInvitation(selection.Text())
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
	if sel.Length() == 0 {
		return false, fmt.Errorf("empty selection")
	}
	file, err := os.Create("/home/elgreco/Desktop/afolder/newfolder/" + strconv.Itoa(num) + ".txt")
	if err != nil {
		return false, fmt.Errorf("unable to create file %s: %s", resp.Request.URL.Path+".txt", err)
	}
	defer file.Close()
	sel.Each(func(i int, selection *goquery.Selection) {
		//_, err = file.WriteString(strings.TrimSpace(selection.Closest("div").Contents().Text()) +
		//	"\n_______________________\n")
		//if err != nil {
		//	log.Printf("unable to write to file %s: %s", resp.Request.URL.Path+".txt", err)
		//}
		text := strings.TrimSpace(selection.Closest("div").Contents().Text())
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "\t", " ")
		text = strings.ReplaceAll(text, "\r", " ")
		resp, err := watson.SendToWatson(text)
		//if err != nil {
		//	return false, fmt.Errorf("error when interacting with watson: %v", err)
		//}
		if err != nil {
			log.Println(err)
		} else {
			_, err = file.WriteString(text)
			if err != nil {
				log.Printf("unable to write to file %s: %s", file.Name(), err)
			}
			_, err = file.WriteString("\n" + resp)
			if err != nil {
				log.Printf("unable to write to file %s: %s", file.Name(), err)
			}
		}
	})
	return true, nil
}
