package main

import "github.com/illfalcon/parser/finder"

func main() {
	// Make HTTP GET request
	//client := &http.Client{
	//	Timeout: 30 * time.Second,
	//}
	//
	//// Make request
	//response, err := client.Get("https://www.devdungeon.com/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer response.Body.Close()
	//
	//// Copy data from the response to standard output
	//n, err := io.Copy(os.Stdout, response.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println("Number of bytes copied to STDOUT:", n)
	_, _ = finder.ContainsDate()
}
