package collyparser

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/illfalcon/parser/pkg/str"

	"github.com/go-shiori/go-readability"

	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

func isWebPage(link string) bool {
	return filepath.Ext(link) == "" || filepath.Ext(link) == ".html" || filepath.Ext(link) == ".php" ||
		filepath.Ext(link) == ".asp"
}

func CrawlAndParse(site string, depth int) error {
	dstTxtFile, _ := os.Create("text2.txt")
	defer dstTxtFile.Close()
	u, err := url.Parse(site)
	if err != nil {
		return errors.Wrapf(err, "error parsing initial url %v", site)
	}
	c := colly.NewCollector(
		colly.AllowedDomains(u.Host),
		colly.MaxDepth(depth),
	)
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})
	if err != nil {
		return errors.Wrap(err, "unable to configure colly parser")
	}
	c.OnRequest(func(r *colly.Request) {
		log.Printf("requested page %v\n", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		log.Printf("received response from %v with code %v\n", r.Request.URL, r.StatusCode)
		err := parse(bytes.NewReader(r.Body), r.Request.URL.String(), dstTxtFile)
		if err != nil {
			log.Printf("error when parsing %s:, %v", r.Request.URL.RawPath, err)
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if isWebPage(e.Request.AbsoluteURL(link)) {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})
	c.Visit(site)
	//c.Wait()
	return nil
}

func parse(body io.Reader, url string, w io.Writer) error {
	article, err := readability.FromReader(body, url)
	if err != nil {
		return errors.Wrapf(err, "unable to parse url %s", url)
	}
	textContent := article.TextContent
	textContent = str.NormalizeNewLines(textContent)
	textContent = str.NormalizePunctuation(textContent)
	scanner := bufio.NewScanner(strings.NewReader(textContent))
	for scanner.Scan() {
		if containsInvitation(scanner.Text()) {
			w.Write([]byte(scanner.Text() + "\n\n\n"))
		}
	}
	return nil
}
