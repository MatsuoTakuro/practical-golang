package ch10

import (
	"log"
	"net/http"
	"time"
)

func middleware() {
	// none()
	logging()
}

func none() {
	http.Handle("/healthz", http.HandlerFunc(Healthz))
	http.ListenAndServe(":8888", nil)
}

func logging() {
	http.Handle("/healthz", MiddlewareLogging(http.HandlerFunc(Healthz)))
	http.ListenAndServe(":8888", nil)
}

func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("start %s\n", r.URL)
		next.ServeHTTP(w, r)
		log.Printf("end %s\n", r.URL)
	})
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK (in Response Body)\n"))
	// panic("panic occurs")
}
