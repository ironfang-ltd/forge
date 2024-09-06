package main

import (
	"net/http"

	"github.com/ironfang-ltd/forge/router"
)

func main() {

	r := router.New()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	err := http.ListenAndServe("127.0.0.1:5000", r)
	if err != nil {
		panic(err)
	}
}
