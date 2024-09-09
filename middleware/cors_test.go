package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors(t *testing.T) {

	tests := []struct {
		name            string
		options         []CorsOption
		method          string
		requestHeaders  map[string]string
		responseHeaders map[string]string
	}{
		{
			name:           "DefaultConfig",
			method:         "GET",
			options:        []CorsOption{},
			requestHeaders: map[string]string{},
			responseHeaders: map[string]string{
				"Vary": "Origin",
			},
		},
		{
			name:    "AnyOrigin",
			options: []CorsOption{},
			method:  "GET",
			requestHeaders: map[string]string{
				"Origin": "http://example.com",
			},
			responseHeaders: map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "*",
			},
		},
		{
			name: "SingleOrigin",
			options: []CorsOption{
				WithAllowedOrigins("http://example.com"),
			},
			method: "GET",
			requestHeaders: map[string]string{
				"Origin": "http://example.com",
			},
			responseHeaders: map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://example.com",
			},
		},
		{
			name: "MultipleOrigins",
			options: []CorsOption{
				WithAllowedOrigins("http://example.com", "http://example.org"),
			},
			method: "GET",
			requestHeaders: map[string]string{
				"Origin": "http://example.org",
			},
			responseHeaders: map[string]string{
				"Vary":                        "Origin",
				"Access-Control-Allow-Origin": "http://example.org",
			},
		},
	}

	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			mw := Cors(tc.options...)

			req, _ := http.NewRequest(tc.method, "http://example.com/endpoint", nil)
			for name, value := range tc.requestHeaders {
				req.Header.Add(name, value)
			}

			res := httptest.NewRecorder()

			mw(final).ServeHTTP(res, req)

			assertResponseHeaders(t, res.Header(), tc.responseHeaders)
		})
	}
}

func assertResponseHeaders(t *testing.T, resHeader http.Header, expected map[string]string) {

	for name, value := range expected {

		got := resHeader.Get(name)

		if got != value {
			t.Errorf("expected header '%s' to be '%s', got '%s'", name, value, got)
		}
	}
}
