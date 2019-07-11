package finder

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/watson-developer-cloud/go-sdk/assistantv2"

	"github.com/illfalcon/parser/internal/db"
	"github.com/illfalcon/parser/internal/watson"
	"github.com/illfalcon/parser/pkg/str"

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

func prepareDoc(htmlText string) string {
	htmlText = strings.ReplaceAll(htmlText, "<em>", "")
	htmlText = strings.ReplaceAll(htmlText, "<strong>", "")
	htmlText = strings.ReplaceAll(htmlText, "<span>", "")
	htmlText = strings.ReplaceAll(htmlText, "</em>", "")
	htmlText = strings.ReplaceAll(htmlText, "</strong>", "")
	htmlText = strings.ReplaceAll(htmlText, "</span>", "")
	reg := regexp.MustCompile(`<!--[^>]*-->`)
	htmlText = reg.ReplaceAllString(htmlText, "")
	return htmlText
}

func WriteDivsWithDate(response *http.Response, dbWriter db.TextWriter) error {
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read body")
	}
	htmlText := string(htmlBytes)
	htmlText = prepareDoc(htmlText)
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
		return fmt.Errorf("empty selection")
	}
	sel.Each(func(i int, selection *goquery.Selection) {
		f := func(s string) error {
			h := str.FindSha1Hash(s)
			resp := &assistantv2.MessageResponse{Output: &assistantv2.MessageOutput{}}
			if contains, _ := dbWriter.ContainsHash(h); !contains {
				if len(s) > 2048 {
					chch := str.ChunkString(s, 2048)
					for _, ch := range chch {
						r, err := watson.Send(ch)
						if err != nil {
							return err
						}
						resp.Output.Intents = append(resp.Output.Intents, r.Output.Intents...)
					}
				} else {
					resp, err = watson.Send(s)
					if err != nil {
						return err
					}
				}
				var intent string
				var confidence float64 = -1
				if resp.Output.Intents != nil && len(resp.Output.Intents) != 0 {
					for _, i := range resp.Output.Intents {
						if *i.Confidence > confidence {
							intent = *i.Intent
							confidence = *i.Confidence
						}
					}
				} else {
					intent = "irrelevant"
					confidence = 0
				}
				err = dbWriter.AddText(s, h, response.Request.URL.String(), intent, confidence)
				if err != nil {
					return err
				}
			}
			return nil
		}
		text := strings.TrimSpace(selection.Closest("div").Contents().Text())
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "\t", " ")
		text = strings.ReplaceAll(text, "\r", " ")
		_ = f(text)
	})
	return nil
}

//func unmarshalWatson(resp string) (map[string]interface{}, error) {
//	var m map[string]interface{}
//	err := json.Unmarshal([]byte(resp), &m)
//	if err != nil {
//		return nil, err
//	}
//	return m, nil
//}
