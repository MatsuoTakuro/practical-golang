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
	go func() {
		err := runDummyServer()
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(5 * time.Millisecond)
	retry()
	retryWithRetryablehttp()
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
	fmt.Println()
	fmt.Println("custom")
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
	log.Println("status:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("body:\n", string(body))
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
		t.logger("%s %s %d %s, duration: %d ms",
			req.Method, req.URL.String(), resp.StatusCode, http.StatusText(resp.StatusCode), time.Since(start).Milliseconds())
	}
	return resp, err
}

func logging() {
	fmt.Println()
	fmt.Println("logging")
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
	log.Println("status:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("body:\n", string(body))
}

type basicAuthRoundTripper struct {
	base     http.RoundTripper
	username string
	password string
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
	log.Println("status:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("body:\n", string(body))
}

func retry() {
	fmt.Println()
	fmt.Println("retry")
	clt := &http.Client{
		Transport: &retryRoundTripper{
			base:     http.DefaultTransport,
			attempts: 3,
			waitTime: 3 * time.Second,
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8111/dummy", nil)
	if err != nil {
		panic(err)
	}

	resp, err := clt.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("finally, status:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("finally, body:", string(body))
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
			log.Println("resp status:", resp.StatusCode)
			return true
		}
	}

	if resp != nil {
		if resp.StatusCode == 429 || (500 <= resp.StatusCode && resp.StatusCode <= 504) {
			log.Println("resp status:", resp.StatusCode)
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
		log.Println("attempt count:", count)
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
	log.Println("attempts are over")

	return resp, err
}

func retryWithRetryablehttp() {
	fmt.Println()
	fmt.Println("retryWithRetryablehttp")
	retryableClt := retryablehttp.NewClient()
	retryableClt.RetryMax = 3
	retryableClt.Backoff = retryablehttp.LinearJitterBackoff // instead of DefaultBackoff
	// get a small amount of jitter centered around one second increasing each retry
	retryableClt.RetryWaitMin = 800 * time.Millisecond
	retryableClt.RetryWaitMax = 1200 * time.Millisecond

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8111/dummy", nil)
	if err != nil {
		panic(err)
	}

	resp, err := retryableClt.StandardClient().Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("finally, status:", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println("finally, body:", string(body))
}
