package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_WithDefaults(t *testing.T) {

	req, err := http.NewRequest("OPTIONS", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	m := Cors()

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("got: `%s`, want: `*`", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("got: `%s`, want: `GET, POST, PUT, DELETE, OPTIONS`", w.Header().Get("Access-Control-Allow-Methods"))
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("got: `%s`, want: `Content-Type, Authorization`", w.Header().Get("Access-Control-Allow-Headers"))
	}

	if w.Header().Get("Access-Control-Expose-Headers") != "Authorization" {
		t.Errorf("got: `%s`, want: `Authorization`", w.Header().Get("Access-Control-Expose-Headers"))
	}

	if w.Header().Get("Access-Control-Max-Age") != "600" {
		t.Errorf("got: `%s`, want: `600`", w.Header().Get("Access-Control-Max-Age"))
	}
}

func TestCors_WithAllowedOrigins(t *testing.T) {

	req, err := http.NewRequest("OPTIONS", "http://example.com/", nil)
	req.Header.Set("Origin", "http://example.com")
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	m := Cors(WithAllowedOrigins("http://example.com"))

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("got: `%s`, want: `http://example.com`", w.Header().Get("Access-Control-Allow-Origin"))
	}
}
