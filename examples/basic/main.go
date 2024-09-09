package main

import (
	"net/http"

	"github.com/ironfang-ltd/go-router"
)

func main() {

	r := router.New()

	// Simple
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// With Param
	r.Get("/:name", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Write([]byte("Hello, " + name + "!"))
	})

	// With Group
	apiGroup := r.Group("/api")

	apiGroup.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World from api!"))
	})

	err := http.ListenAndServe("127.0.0.1:5000", r)
	if err != nil {
		panic(err)
	}
}
