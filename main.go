package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/illfalcon/parser/finder"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	fileName := "/home/elgreco/Desktop/afolder/newfile.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error when opening file %s\n", fileName)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var num int
	for scanner.Scan() {
		url := scanner.Text()
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
		_, err = finder.WriteDivsWithDate(resp, num)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(num)
		num++
		resp.Body.Close()
	}

}
