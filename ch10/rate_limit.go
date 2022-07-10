package ch10

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func rateLimit() {
	http.Handle("/test", LimitHandler(http.HandlerFunc(hello)))
	http.ListenAndServe(":18888", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world!")
}

var limiter = rate.NewLimiter(rate.Every(time.Second/1), 1)

func LimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
