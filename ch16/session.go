package ch16

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func session() {
	// tokenWithContext()
	// tokenWithHttpRequest()
	loggerWithContext()
}

var tokenContextKey = struct{}{} // make it private

func tokenWithContext() {
	ctx := context.Background()
	token := "xxxx"
	ctx2 := RegisterToken(ctx, token)
	t, _ := RetrieveToken(ctx2)
	fmt.Println("token", t)
}

// dedicate function for registering ‘token'
func RegisterToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

// dedicate function for retrieving ‘token'
func RetrieveToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok {
		return "", errors.New("TOKEN IS NOT REGISTERED")
	}
	return token, nil
}

func tokenWithHttpRequest() {
	http.HandleFunc("/token", func(writer http.ResponseWriter, req *http.Request) {
		token := NewToken(10)
		req2 := RegisterToken2(req, token)
		t, _ := RetrieveToken2(req2)
		fmt.Fprintf(writer, "token: %s\n", t)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NewToken(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

// v2; dedicate function for registering ‘token'
func RegisterToken2(req *http.Request, token string) *http.Request {
	ctx := context.WithValue(req.Context(), tokenContextKey, token)
	return req.WithContext(ctx)
}

// v2; dedicate function for retrieving ‘token'
func RetrieveToken2(req *http.Request) (string, error) {
	token, ok := req.Context().Value(tokenContextKey).(string)
	if !ok {
		return "", errors.New("TOKEN IS NOT REGISTERED")
	}
	return token, nil
}

var logContextKey = struct{}{}

func loggerWithContext() {
	http.HandleFunc("/logger", func(writer http.ResponseWriter, req *http.Request) {
		req2 := StartLogging(req)
		err := CountAccess(req2.Context())
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("Standard error message."))
		}
		l, err := FinishLogging(req2)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("Standard error message."))
		}
		fmt.Fprintf(writer, "LogRecord; start: %v Duration: %v DBAccessCount: %d\n", l.start, l.Duration, l.DBAccessCount)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type LogRecord struct {
	start         time.Time
	Duration      time.Duration
	DBAccessCount int
}

func StartLogging(req *http.Request) *http.Request {
	l := &LogRecord{
		start:         time.Now(),
		Duration:      0,
		DBAccessCount: 0,
	}
	ctx := context.WithValue(req.Context(), logContextKey, l)
	return req.WithContext(ctx)
}

func CountAccess(ctx context.Context) error {
	l, ok := ctx.Value(logContextKey).(*LogRecord)
	if !ok {
		return errors.New("LOGGER IS NOT REGISTERED")
	}
	l.DBAccessCount += 1
	return nil
}

func FinishLogging(req *http.Request) (*LogRecord, error) {
	l, ok := req.Context().Value(logContextKey).(*LogRecord)
	if !ok {
		return nil, errors.New("LOGGER IS NOT REGISTERED")
	}
	l.Duration = time.Since(l.start)
	return l, nil
}
