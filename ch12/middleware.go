package ch12

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

// curl localhost:8000/
func middleware() {
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	mux := http.NewServeMux()
	mux.Handle("/", WithLogger(logger)(http.HandlerFunc(HelloHandler)))
	fmt.Println("start listening at :8000")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

type contextKey string

const logKey contextKey = "log"

func WithLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceId := r.Header.Get("Trace-Id")
			if traceId == "" {
				traceId = xid.New().String() // generate globally unique ID
			}
			logger = logger.With().Str("trace-id", traceId).Logger()
			ctx := context.WithValue(r.Context(), logKey, logger)

			writer := &responseWriterWrapper{
				writer:  w,
				request: r,
				start:   time.Now(),
			}
			next.ServeHTTP(writer, r.WithContext(ctx))
			logger.Info().Object("httpRequest", writer).Send()
		})
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "hello world\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// ctx := r.Context()
	// logger := GetLogger(ctx)
	// writer := &responseWriterWrapper{
	// 	writer:  w,
	// 	request: r,
	// 	start:   time.Now(),
	// }
	// logger.Info().Object("httpRequest", writer).Send()
}

func GetLogger(ctx context.Context) zerolog.Logger {
	return ctx.Value(logKey).(zerolog.Logger)
}

type responseWriterWrapper struct {
	status       int
	responseSize int
	writer       http.ResponseWriter
	request      *http.Request
	start        time.Time
}

var _ http.Flusher = &responseWriterWrapper{}

func (r *responseWriterWrapper) Flush() {
	flusher := r.writer.(http.Flusher)
	flusher.Flush()
	r.status = 200
}

var _ http.ResponseWriter = &responseWriterWrapper{}

// impl ResponseWriter IF
func (r *responseWriterWrapper) Header() http.Header {
	return r.writer.Header()
}

// impl ResponseWriter IF
func (r *responseWriterWrapper) Write(content []byte) (int, error) {
	r.responseSize += len(content)
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.writer.Write(content)
}

// impl ResponseWriter IF
func (r *responseWriterWrapper) WriteHeader(statusCode int) {
	r.status = statusCode
	r.writer.WriteHeader(statusCode)
}

var _ zerolog.LogObjectMarshaler = &responseWriterWrapper{}

// impl LogObjectMarshaler IF
func (r *responseWriterWrapper) MarshalZerologObject(e *zerolog.Event) {
	e.Str("requestMethod", r.request.Method)
	e.Str("requestUrl", r.request.URL.String())
	e.Int64("requestSize", r.request.ContentLength)
	e.Int("status", r.status)
	e.Int("responseSize", r.responseSize)
	e.Str("referer", r.request.Header.Get("Referer"))
	e.Str("latency", time.Since(r.start).String())
	e.Bool("cacheHit", r.status == 304)
	forwarded := r.request.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		e.Str("remoteIp", forwarded)
	} else {
		e.Str("remoteIp", r.request.RemoteAddr)
	}
	e.Str("protocol", r.request.Proto)
}
