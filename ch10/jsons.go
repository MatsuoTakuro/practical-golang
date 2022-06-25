package ch10

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Comment struct {
	Message  string `json:"message,omitempty" validate:"required,min=1,max=140"`
	UserName string `json:"user_name,omitempty" validate:"required,min=1,max=15"`
}

var val = validator.New()
var ve validator.ValidationErrors

// curl  localhost:8888/comments --data '{"message": "test", "user_name": "gotaro"}'
// curl  localhost:8888/comments --data '{"message": "test2", "user_name": "gotaro2"}'
// curl  localhost:8888/comments
// curl  localhost:8888/comments --data '{"message": "test", "user_name": "1234567890123456"}'
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
			if err := val.Struct(c); err != nil {
				var out []string
				if errors.As(err, &ve) {
					for _, fe := range ve {
						switch fe.Field() {
						case "Message": // struct field name, not json tag name
							out = append(out, "message is from 1 to 140 chars")
						case "UserName":
							out = append(out, "user_name is from 1 to 15 chars")
						}
					}
				}
				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, strings.Join(out, ",")), http.StatusBadRequest)
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
