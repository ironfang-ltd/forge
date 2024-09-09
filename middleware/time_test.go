package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTime(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	m := Time()

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("X-Request-Time-Ms") == "" {
		t.Error("X-Request-Time-Ms header is not set")
	}
}
