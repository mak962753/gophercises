package main

// yaml parse
// go get gopkg.in/yaml.v2

import (
	"flag"
	"fmt"
	urlshort "gophercises/02-urlshort"
	"io/ioutil"
	"net/http"
)

func main() {
	yamlFile := flag.String("yaml", "", "if set specifies YAML file that contain url mappings.")
	jsonFile := flag.String("json", "", "if set specifies YAML file that contain url mappings.")

	flag.Parse()

	mux := defaultMux()

	handler := addMapHander(mux)
	handler = addYamlHandler(yamlFile, handler)
	handler = addJsonHandler(jsonFile, handler)

	fmt.Println("Starting HTTP server on :8080... ")

	_ = http.ListenAndServe(":8080", handler)
}

func addMapHander(fallback http.Handler) http.HandlerFunc {
	paths := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	return urlshort.MapHandler(paths, fallback)
}
func addYamlHandler(file *string, fallback http.HandlerFunc) http.HandlerFunc {
	const yamlUrls = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	data := []byte(yamlUrls)
	if file != nil {
		var err error
		data, err = ioutil.ReadFile(*file)
		if err != nil {
			fmt.Println("addYamlHandler: Failed to read yaml file", err)
			return fallback
		}
	}
	yamlHandler, yamlErr := urlshort.YAMLHandler(data, fallback)
	if yamlErr == nil {
		return yamlHandler
	}
	fmt.Println("addYamlHandler: to initialize YAML map handler.", yamlErr)
	return fallback
}

func addJsonHandler(file *string, fallback http.HandlerFunc) http.HandlerFunc {
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println("addJsonHandler: Failed to read json file", err)
		return fallback
	}
	jsonHandler, err := urlshort.JSONHandler(data, fallback)
	if err != nil {
		fmt.Println("addJsonHandler: Failed to initialize json handler", err)
		return fallback
	}
	return jsonHandler
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello!!!!")
	})
	return mux
}
