package ch11

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func normal() {
	// httpGet()
	// httpPost()
	httpClient()
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

type User struct {
	Name string
	Addr string
}

func httpPost() {
	u := User{
		Name: "O'Reilly Japan",
		Addr: "東京都新宿区四谷坂町",
	}

	payload, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://example.com", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func httpClient() {
	clt := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer XXX...XXX")

	resp, err := clt.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status code is\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
