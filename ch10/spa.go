package ch10

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
)

//go:embed vite-project/dist/*
var assets embed.FS

func tryRead(requestedPath string, w http.ResponseWriter) error {
	log.Println(requestedPath)

	f, err := assets.Open(path.Join("vite-project/dist", requestedPath))
	if err != nil {
		return err
	}
	defer f.Close()

	stat, _ := f.Stat()
	if stat.IsDir() {
		return errors.New("path is dir")
	}

	ext := filepath.Ext(requestedPath)
	var contentType string
	if m := mime.TypeByExtension(ext); m != "" {
		contentType = m
	} else {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	return nil
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := tryRead(r.URL.Path, w)
	if err == nil {
		return
	}

	err = tryRead("index.html", w)
	if err != nil {
		panic(err)
	}
}

func newHandler() http.Handler {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		// curl localhost:8000/api/test
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test\n"))
		})
	})
	// curl localhost:8000/
	// curl localhost:8000/api
	// curl localhost:8000/hello
	router.NotFound(notFoundHandler)

	return router
}

func spa() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	server := &http.Server{
		Addr:    ":8000",
		Handler: newHandler(),
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()
	fmt.Println("start receiving at :8000")
	fmt.Println(os.Stderr, server.ListenAndServe())
}
