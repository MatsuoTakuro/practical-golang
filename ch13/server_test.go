package ch13

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerRequest(t *testing.T) {
	s := httptest.NewServer(initServer())
	r, err := http.Get(s.URL + "/fortune")
	if err != nil {
		t.Errorf("http get err should be nil: %v", err)
	}
	defer r.Body.Close()
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Errorf("json decode err should be nil: %v", err)
		return
	}
	if body["fortune"] != "Daikichi" {
		t.Errorf("result should be Daikichi, but %s", body["fortune"])
	}
}

func TestHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/fortune", nil)
	w := httptest.NewRecorder()
	fortuneHandler(w, r)
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Errorf("json unmarshal err should be nil: %v", err)
		return
	}
	if body["fortune"] != "Daikichi" {
		t.Errorf("result should be Daikichi, but %s", body["fortune"])
	}
}
