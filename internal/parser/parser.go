package parser

import (
	"log"
	"net/http"
	"time"

	"github.com/illfalcon/parser/internal/db"
)

func Parse() {
	//db.Prepare()
	service := db.CreateSqliteService()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	urls, err := service.GetUnparsedURLs()
	if err != nil {
		log.Fatal(err)
	}
	for _, url := range urls {
		//url := scanner.Text()
		// Make request
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("error when getting url %s\n", url)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Printf("error when getting url %s, status: %s\n", url, resp.Status)
		}
		err = writeDivsWithDate(resp, &service)
		if err != nil {
			log.Print(err)
		}
		err = service.SetURLParsed(url)
		if err != nil {
			log.Print(err)
		}
		resp.Body.Close()
	}

}
