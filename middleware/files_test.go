package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFiles_WithDefaults(t *testing.T) {

	req, err := http.NewRequest("GET", "/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("filePath", "test.txt")
	w := httptest.NewRecorder()

	m := Files(WithDirectory("../web/static"))

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200, got:", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("response content type is not text/plain")
	}

	if w.Body.String() != "This is just a test file for testing the files middleware." {
		t.Error("response body is not 'This is just a test file for testing the files middleware.', got:", w.Body.String())
	}

}

func TestFiles_WithPathOutsideRoot(t *testing.T) {

	req, err := http.NewRequest("GET", "/../go.mod", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.SetPathValue("filePath", "../go.mod")
	w := httptest.NewRecorder()

	m := Files(WithDirectory("../web/static"))

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("response code is not 404, got:", w.Code)
	}
}
