package ch10

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Comment struct {
	Message  string `json:"message,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

// curl  localhost:8888/comments --data '{"message": "test", "user_name": "gotaro"}'
// curl  localhost:8888/comments --data '{"message": "test2", "user_name": "gotaro2"}'
// curl  localhost:8888/comments
func jsons() {
	var mutex = &sync.RWMutex{}
	comments := make([]Comment, 0, 100)

	http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		switch r.Method {
		case http.MethodGet:
			mutex.RLock()

			if err := json.NewEncoder(w).Encode(comments); err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
				return
			}
			mutex.RUnlock()

		case http.MethodPost:
			var c Comment
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
				return
			}
			mutex.Lock()
			comments = append(comments, c)
			mutex.Unlock()
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write([]byte(`"status": "created"`))
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
			}

		default:
			http.Error(w, `{"status":"permits only GET or POST}`, http.StatusMethodNotAllowed)
		}
	})
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal(err)
	}
}
