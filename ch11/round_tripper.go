package ch11

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func roundTripper() {
	// custom()
	// logging()
	basicAuth()
}

type customRoundTripper struct {
	base http.RoundTripper
}

func (c customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// pre-processing
	log.Println("pre-processing")

	resp, err := c.base.RoundTrip(req)

	// post-processing
	log.Println("post-processing")

	return resp, err
}

func custom() {
	clt := &http.Client{
		Transport: &customRoundTripper{
			base: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}

	resp, err := clt.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status code is\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Body is as follows;\n", string(body))
}

type loggingRoundTripper struct {
	base   http.RoundTripper
	logger func(string, ...any)
}

func (t *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.logger == nil {
		t.logger = log.Printf
	}

	start := time.Now()
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		t.logger("%s %s %d %s, duration: %d",
			req.Method, req.URL.String(), resp.StatusCode, http.StatusText(resp.StatusCode), time.Since(start).Milliseconds())
	}
	return resp, err
}

func logging() {
	clt := &http.Client{
		Transport: &loggingRoundTripper{
			base: http.DefaultTransport,
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}

	resp, err := clt.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status code is\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Body is as follows;\n", string(body))
}

type basicAuthRoundTripper struct {
	username string
	password string
	base     http.RoundTripper
}

func (rt *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(rt.username, rt.password)
	log.Println("Authorization:", req.Header["Authorization"])
	return rt.base.RoundTrip(req)
}

type User2 struct {
	username string
	password string
}

func basicAuth() {
	u := &User2{
		username: "MatsuoTakuro",
		password: "1234567890",
	}

	clt := &http.Client{
		Transport: &basicAuthRoundTripper{
			username: u.username,
			password: u.password,
			base:     http.DefaultTransport,
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}

	resp, err := clt.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Status code is\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Body is as follows;\n", string(body))
}
