package ch13

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func server() {
	log.Println("open at :8080")
	log.Println(http.ListenAndServe(":8080", initServer()))
}

// curl -v localhost:8080/fortune
func initServer() http.Handler {
	r := chi.NewRouter()
	r.Get("/fortune", fortuneHandler)
	return r
}

func fortuneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, err := io.WriteString(w, `{"fortune": "Daikichi"}`)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
	}
}
