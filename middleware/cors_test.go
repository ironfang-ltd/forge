package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_OnlyHandleOPTIONS(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	m := Cors()

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("response code is: %d, expected: %d", w.Code, http.StatusOK)
	}

	if w.Header().Get("Vary") != "" {
		t.Errorf("vary header is '%s', expected: ''", w.Header().Get("Vary"))
	}
}

func TestCors_WithDefaults(t *testing.T) {

	req, err := http.NewRequest("OPTIONS", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	m := Cors()

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("response code is: %d, expected: %d", w.Code, http.StatusOK)
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("access-control-allow-origin header is '%s', expected: %s", w.Header().Get("Access-Control-Allow-Origin"), "*")
	}

	if w.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("access-control-allow-methods header is '%s', expected: %s", w.Header().Get("Access-Control-Allow-Methods"), "GET, POST, PUT, DELETE, OPTIONS")
	}

	if w.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("access-control-allow-headers header is '%s', expected: %s", w.Header().Get("Access-Control-Allow-Headers"), "Content-Type, Authorization")
	}

	if w.Header().Get("Access-Control-Expose-Headers") != "Authorization" {
		t.Errorf("access-control-expose-headers header is '%s', expected: %s", w.Header().Get("Access-Control-Expose-Headers"), "Authorization")
	}

	if w.Header().Get("Access-Control-Max-Age") != "600" {
		t.Errorf("access-control-max-age header is '%s', expected: %s", w.Header().Get("Access-Control-Max-Age"), "600")
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

func TestCors_WithBadOrigin(t *testing.T) {

	req, err := http.NewRequest("OPTIONS", "http://example.com/", nil)
	req.Header.Set("Origin", "http://test.com")
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
