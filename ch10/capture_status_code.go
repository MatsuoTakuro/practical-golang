package ch10

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	log.Printf("Response Body: %s", b)
	return lrw.ResponseWriter.Write(b)
}

func wrappedHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode

		log.Printf("%d %s", statusCode, http.StatusText(statusCode))
	})
}

func captureStatusCode() {
	http.Handle("/healthz", wrappedHandlerWithLogging(http.HandlerFunc(Healthz)))
	http.ListenAndServe(":8888", nil)
}
