package ch10

import (
	"net/http"
	"time"
)

func timeout() {
	h := MiddlewareLogging(http.HandlerFunc(Healthz))
	http.Handle("/healthz", http.TimeoutHandler(h, 3*time.Second, "request timeout\n"))
	http.ListenAndServe(":8888", nil)
}
