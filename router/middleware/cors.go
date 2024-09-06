package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ironfang-ltd/forge/router"
)

type CorsOption func(*CorsOptions)

type CorsOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposedHeaders []string
	MaxAge         int
}

func WithAllowedOrigins(origins ...string) func(*CorsOptions) {
	return func(opts *CorsOptions) {
		opts.AllowedOrigins = origins
	}
}

func Cors(options ...CorsOption) router.Middleware {

	opts := &CorsOptions{
		AllowedOrigins: []string{},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		ExposedHeaders: []string{"Authorization"},
		MaxAge:         600,
	}

	for _, option := range options {
		option(opts)
	}

	methods := strings.Join(opts.AllowedMethods, ", ")
	headers := strings.Join(opts.AllowedHeaders, ", ")
	exposed := strings.Join(opts.ExposedHeaders, ", ")

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodOptions {
				next(w, r)
				return
			}

			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			origin := r.Header.Get("Origin")

			if origin == "" {
				next(w, r)
				return
			}

			if len(opts.AllowedOrigins) > 0 {
				for _, allowedOrigin := range opts.AllowedOrigins {
					if origin == allowedOrigin {
						w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
						break
					}
				}
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Expose-Headers", exposed)
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(opts.MaxAge))

			w.WriteHeader(http.StatusOK)
		}
	}
}
