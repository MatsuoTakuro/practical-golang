package ch10

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func multiplexer() {
	// normal()
	withChi()
}

func normal() {
	yes := 0
	no := 0
	mux := http.NewServeMux()

	// curl localhost:8080/asset/ch1/
	mux.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir("."))))

	// curl localhost:8080/poll/y
	mux.HandleFunc("/poll/y", func(w http.ResponseWriter, r *http.Request) {
		yes++
	})

	// curl localhost:8080/poll/n
	mux.HandleFunc("/poll/n", func(w http.ResponseWriter, r *http.Request) {
		no++
	})

	// curl localhost:8080/result
	mux.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "for: %d, against: %d\n", yes, no)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func withChi() {
	yes := 0
	no := 0
	r := chi.NewRouter()

	r.Post("/poll/{answer}", func(w http.ResponseWriter, r *http.Request) {
		switch chi.URLParam(r, "answer") {
		// curl -X POST localhost:8080/poll/y
		case "y":
			yes++
		// curl -X POST localhost:8080/poll/n
		case "n":
			no++
		// curl -X POST localhost:8080/poll/1
		default:
			http.Error(w, fmt.Sprintf(`{"status":"%d"}`, http.StatusBadRequest), http.StatusBadRequest)
		}
	})

	// curl localhost:8080/result
	r.Get("/result", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "for: %d, against: %d\n", yes, no)
	})

	// curl localhost:8080/asset/ch1/
	r.Handle("/asset/*", http.StripPrefix("/asset/", http.FileServer(http.Dir("."))))

	log.Fatal(http.ListenAndServe(":8080", r))
}
