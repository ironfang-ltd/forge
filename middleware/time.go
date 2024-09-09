package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ironfang-ltd/go-router"
)

type customWriter struct {
	w      http.ResponseWriter
	before func()
}

func (cw customWriter) WriteHeader(code int) {
	if cw.before != nil {
		cw.before()
	}
	cw.w.WriteHeader(code)
}

func (cw customWriter) Write(b []byte) (int, error) {
	if cw.before != nil {
		cw.before()
	}
	return cw.w.Write(b)
}

func (cw customWriter) Header() http.Header {
	return cw.w.Header()
}

func Time() router.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			writtenHeader := false

			cw := customWriter{
				w: w,
				before: func() {
					if writtenHeader {
						return
					}
					writtenHeader = true

					taken := time.Since(start)
					w.Header().Set("X-Request-Time-Ms", strconv.FormatInt(taken.Milliseconds(), 10))
				},
			}

			next(cw, r)
		}
	}
}
