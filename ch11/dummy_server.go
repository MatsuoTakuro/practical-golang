package ch11

import (
	"fmt"
	"log"
	"net/http"
)

type dummyHandler struct {
	count int
}

const (
	MAX_FAILIRES = 5
)

var reqCount int

func runDummyServer() error {
	fmt.Println()
	fmt.Println("runDummyServer")

	mux := http.NewServeMux()
	mux.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
		if reqCount <= MAX_FAILIRES {
			w.WriteHeader(http.StatusTooManyRequests)
			reqCount += 1
			if _, err := fmt.Fprintf(w, "%s",
				fmt.Sprintf("your request fails. please retry. (dummy handler's count: %d)", reqCount)); err != nil {
				log.Fatal(err)
			}
		} else {
			if _, err := fmt.Fprintf(w, "%s",
				fmt.Sprintf("your request is successful! (dummy handler's count: %d)", reqCount)); err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("dummy handler's count: %d", reqCount)
	})

	port := ":8111"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("start listening at %s (request consecutive dummy failures = %d)", port, MAX_FAILIRES)
	return server.ListenAndServe()
}
