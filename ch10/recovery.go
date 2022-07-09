package ch10

import (
	"log"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func recoveryFromPanic() {
	http.Handle("/healthz", Recovery(wrappedHandlerWithLogging(http.HandlerFunc(Healthz))))
	http.ListenAndServe(":8888", nil)
}
