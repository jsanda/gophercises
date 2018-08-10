package exercises

import (
	"net/http"
	"fmt"
)

type urlShortener struct {}

func NewUrlShortener() *urlShortener {
	return &urlShortener{}
}

func (*urlShortener) Run() error {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the YAMLHandler as the fallback
	json := `
[
  {
    "path": "/json-godoc",
    "url": "https://golang.org/pkg/encoding/json/"
  },
  {
    "path": "/urlshort-readme",
    "url": "https://github.com/gophercises/urlshort/blob/master/README.md"
  }
]
`
	jsonHandler, err := JSONHandler([]byte(json), yamlHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)


	return nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
