package exercises

import (
	"net/http"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return redirect(pathsToUrls, fallback)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	return redirect(paths, fallback), nil
}

func parseYAML(yml []byte) (map[string]string, error) {
	mappings := make([]map[string]string, 0)
	if err := yaml.Unmarshal(yml, &mappings); err != nil {
		return nil, err
	}
	paths := make(map[string]string)
	for _, v := range mappings {
		paths[v["path"]] = v["url"]
	}

	return paths, nil
}

func redirect(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if val, ok := paths[req.URL.Path]; ok {
			http.Redirect(w, req, val, 300)
		} else {
			fallback.ServeHTTP(w, req)
		}
	}
}
