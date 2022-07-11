package ch11

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func normal() {
	httpGet()
}

func httpGet() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Not OK", resp.Status)
	}
	fmt.Println("Status code is\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Body is as follows;\n", string(body))
}
