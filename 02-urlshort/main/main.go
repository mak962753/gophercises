package main

// yaml parse
// go get gopkg.in/yaml.v2

import (
	"fmt"
	urlshort "gophercises/02-urlshort"
	"net/http"
)

func main() {
	fmt.Println("test")

	paths := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mux := defaultMux()

	handler := urlshort.MapHandler(paths, mux)
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	if yamlHandler, yamlErr := urlshort.YAMLHandler([]byte(yaml), handler); yamlErr == nil {
		handler = yamlHandler
	} else {
		fmt.Println("Failed to initialize YAML map handler.", yamlErr)
	}

	fmt.Println("Starting HTTP server on :8080... ")

	_ = http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello!!!!")
	})
	return mux
}
