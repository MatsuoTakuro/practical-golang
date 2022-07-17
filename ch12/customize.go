package ch12

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

func customize() {
	// withZerolog()
	withUberZap()
}

func withZerolog() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("panic")
	}

	http.HandleFunc("/test", handler)

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	server := &http.Server{
		Addr:     ":18888",
		ErrorLog: log.New(logger, "", 0),
	}

	logger.Fatal().Msgf("server: %v", server.ListenAndServe())
}

type logForwarder struct {
	l *zap.SugaredLogger
}

func (lfw *logForwarder) Write(p []byte) (int, error) {
	lfw.l.Errorw(string(p))
	return len(p), nil
}

func withUberZap() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		panic("panic2")
	}

	http.HandleFunc("/test2", handler)

	l, err := zap.NewDevelopment()
	if err != nil {
		l = zap.NewNop()
	}
	logger := l.Sugar()

	server := &http.Server{
		Addr:     ":18888",
		ErrorLog: log.New(&logForwarder{l: logger}, "", 0),
	}

	logger.Fatal("server: %v", server.ListenAndServe())
}
