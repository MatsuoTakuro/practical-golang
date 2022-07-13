package ch11

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func roundTripper() {
	// custom()
	// logging()
	// basicAuth()
	// retry()
	retryWithJitter()
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

func retry() {
	clt := &http.Client{
		Transport: &retryRoundTripper{
			base:     http.DefaultTransport,
			attempts: 5,
			waitTime: 5 * time.Second,
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

type retryRoundTripper struct {
	base     http.RoundTripper
	attempts int
	waitTime time.Duration
}

func (rt *retryRoundTripper) shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}
	}

	if resp != nil {
		if resp.StatusCode == 429 || (500 <= resp.StatusCode && resp.StatusCode <= 504) {
			return true
		}
	}

	return false
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	for count := 0; count < rt.attempts; count++ {
		resp, err = rt.base.RoundTrip(req)

		if !rt.shouldRetry(resp, err) {
			return resp, err
		}

		// donot use time.Sleep
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(rt.waitTime): // wait for retry
		}
	}

	return resp, err
}

func retryWithJitter() {
	clt := retryablehttp.NewClient()
	clt.RetryMax = 2

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/", nil)
	if err != nil {
		panic(err)
	}

	resp, err := clt.StandardClient().Do(req)
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
