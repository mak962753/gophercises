package urlshort

import (
	"encoding/json"
	"fmt"
	yaml2 "gopkg.in/yaml.v2"
	"net/http"
)

func MapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch url, ok := paths[r.RequestURI]; ok {
		case true:
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusFound)
		default:
			fallback.ServeHTTP(w, r)
		}
	}
}
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsed []entry
	err := json.Unmarshal(jsonData, &parsed)
	fmt.Println("JSONHandler(): initialized", parsed)
	if err != nil {
		return nil, err
	}
	paths := map[string]string{}
	for _, e := range parsed {
		paths[e.Path] = e.Url
	}
	return MapHandler(paths, fallback), nil
}

// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
type entry struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsed []entry
	err := yaml2.Unmarshal(yaml, &parsed)
	if err != nil {
		return nil, err
	}
	fmt.Println("YAMLHandler(): initialized", parsed)
	paths := map[string]string{}
	for _, e := range parsed {
		paths[e.Path] = e.Url
	}

	return MapHandler(paths, fallback), nil
}
