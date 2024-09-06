package router

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Get(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}
}

func TestRouter_GetMux(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}
}

func TestRouter_GetWithParam(t *testing.T) {

	req, _ := http.NewRequest("GET", "/test-value", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Get("/:param", func(w http.ResponseWriter, r *http.Request) {

		param := r.PathValue("param")

		_, _ = w.Write([]byte(param))
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Body.String() != "test-value" {
		t.Error("response body is not 'test-value'")
	}
}

func TestRouter_GetWithMiddleware(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			next(w, r)
		}
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("X-Test") != "test" {
		t.Error("response header X-Test is not 'test'")
	}
}

func TestRouter_GetWithGroup(t *testing.T) {

	req, _ := http.NewRequest("GET", "/group/endpoint", nil)
	w := httptest.NewRecorder()

	r := New()

	g := r.Group("/group")

	g.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			next(w, r)
		}
	})

	g.Get("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("group-endpoint"))
	})

	r.ServeHTTP(w, req)

	for _, r := range r.GetRoutes() {
		fmt.Printf("route: %s %s\n", r.Method, r.Path)
	}

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("X-Test") != "test" {
		t.Error("response header X-Test is not 'test'")
	}
}

func TestRouter_OptionsWithMultipleMiddleware(t *testing.T) {

	req, _ := http.NewRequest("OPTIONS", "/group/endpoint", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			next(w, r)
		}
	})

	g := r.Group("/group")

	g.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-2", "test-2")
			next(w, r)
		}
	})

	g.Get("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("group-endpoint"))
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("response code is not 404")
	}

	if w.Header().Get("X-Test") != "test" {
		t.Error("response header X-Test is not 'test'")
	}

	if w.Header().Get("X-Test-2") != "test-2" {
		t.Error("response header X-Test-2 is not 'test-2'")
	}
}

func TestRouter_Static(t *testing.T) {

	req, _ := http.NewRequest("GET", "/test.txt", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Static("/", "./testdata")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("response code is not 200")
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("response content type is not text/plain")
	}

	if w.Body.String() != "hello world" {
		t.Error("response body is not 'hello world'")
	}
}

func TestRouter_StaticWithPath(t *testing.T) {

	req, _ := http.NewRequest("GET", "/files/test.txt", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Static("/files", "./testdata")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("response code is not 200, got: %d", w.Code)
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Error("response content type is not text/plain")
	}

	if w.Body.String() != "hello world" {
		t.Error("response body is not 'hello world'")
	}
}

func TestRouter_NodeOrder(t *testing.T) {

	r := New()

	r.Static("/files", "./testdata")

	r.Get("/:param", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {})

	routes := r.GetRoutes()

	if routes[0].Path != "/test" {
		t.Error("first route is not '/test', but ", routes[0].Path)
	}

	if routes[1].Path != "/:param" {
		t.Error("second route is not '/:param', but ", routes[1].Path)
	}

	if routes[2].Path != "/files*" {
		t.Error("third route is not '/files*', but ", routes[2].Path)
	}
}

func BenchmarkFib10(b *testing.B) {

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r := New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		r.ServeHTTP(w, req)
	}
}
