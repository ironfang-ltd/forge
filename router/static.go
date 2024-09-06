package router

import (
	"net/http"
	"strings"
)

func StaticFileHandler(router *router, path, dir string) http.HandlerFunc {

	fs := http.Dir(dir)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {

			if router.config.MethodNotAllowedHandler != nil {
				router.config.MethodNotAllowedHandler(w, r)
				return
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		filePath := strings.TrimPrefix(r.URL.Path, path)

		f, err := fs.Open(filePath)
		if err != nil {

			if router.config.NotFoundHandler != nil {
				router.config.NotFoundHandler(w, r)
				return
			}

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
			if router.config.NotFoundHandler != nil {
				router.config.NotFoundHandler(w, r)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	}
}
