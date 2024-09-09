package middleware

import (
	"fmt"
	"net/http"

	"github.com/ironfang-ltd/go-router"
)

type FilesOption func(*FilesOptions)

type FilesOptions struct {
	Directory string
}

func WithDirectory(dir string) func(*FilesOptions) {
	return func(opts *FilesOptions) {
		opts.Directory = dir
	}
}

func Files(options ...FilesOption) router.Middleware {

	opts := &FilesOptions{
		Directory: "./web/static",
	}

	for _, option := range options {
		option(opts)
	}

	fs := http.Dir(opts.Directory)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodGet {

				filePath := r.PathValue("filePath")
				if filePath == "" {
					next(w, r)
					return
				}

				f, err := fs.Open(filePath)
				if err != nil {
					fmt.Println("Error opening file: " + err.Error())
					w.WriteHeader(http.StatusNotFound)
					return
				}
				defer f.Close()

				fi, err := f.Stat()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if fi.IsDir() {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)

				return
			}

			next(w, r)
		}
	}
}
