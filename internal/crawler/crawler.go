package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/illfalcon/parser/pkg/str"

	"github.com/illfalcon/parser/internal/db"
)

const depth int = 2

type tuple struct {
	list  []string
	level int
}

var tokens = make(chan struct{}, 20)

func crawlWithDeepness(src, host string, service db.URLWriter, level, depth int) tuple {
	if level < depth {
		return tuple{crawl(src, host, service), level + 1}
	}
	return tuple{nil, level + 1}
}

func crawl(src, host string, service db.URLWriter) []string {
	h, err := getPageAndFindHash(src)
	if err != nil {
		log.Println(err)
	}
	if b, _ := service.ContainsURL(src); !b {
		err := service.AddURL(src, h)
		if err != nil {
			log.Println(err)
			return nil
		}
		tokens <- struct{}{}
		list, err := exrtactNonDocsWithinOrigin(src, host)
		<-tokens // release the token
		if err != nil {
			log.Print(err)
			return nil
		}
		return list
	} else {
		oldHash, err := service.GetURLHash(src)
		if err != nil {
			log.Println(err)
			return nil
		}
		if oldHash != h {
			tokens <- struct{}{}
			list, err := exrtactNonDocsWithinOrigin(src, host)
			<-tokens // release the token
			if err != nil {
				log.Print(err)
				return nil
			}
			_ = service.SetURLUnparsed(src)
			return list
		}
	}
	return nil
}

func Crawl() {
	service := db.CreateSqliteService()
	worklist := make(chan tuple, 20)
	//if err != nil {
	//	log.Printf("error creating file: %s\n", err)
	//	write = func(src string) error {
	//		fmt.Println(src)
	//		return nil
	//	}
	//} else {
	//	write = func(src string) error {
	//		_, err = f.WriteString(src + "\n")
	//		if err != nil {
	//			return err
	//		}
	//		return nil
	//	}
	//}
	//defer f.Close()
	urls, err := service.GetAllLandings()
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range urls {
		newHash, err := getPageAndFindHash(u)
		if err != nil {

		}
		oldHash, err := service.GetURLHash(u)
		if err != nil {
			log.Println(err)
		}
		if newHash == oldHash {
			continue
		}
		src, err := url.Parse(u)
		if err != nil {
			fmt.Println("cannot resolve origin")
			continue
		}
		var n int
		n++
		go func() { worklist <- tuple{[]string{u}, 0} }()
		seen := make(map[string]bool)
		for ; n > 0; n-- {
			list := <-worklist
			for _, link := range list.list {
				if !seen[link] {
					seen[link] = true
					n++
					go func(link string) {
						worklist <- crawlWithDeepness(link, src.Host, &service, list.level, depth)
					}(link)
				}
			}
		}
	}
}

func getPageAndFindHash(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "", fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	hash := str.FindSha1Hash(string(bb))
	return hash, nil
}
