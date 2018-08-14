package urlshort

import (
	"net/http"
	"fmt"
	"github.com/coreos/bbolt"
	"log"
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

	db, err := openDB()
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}
	if err = initDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %s", err)
	}
	boltHandler := BoltHandler(db, jsonHandler)


	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)

	return nil
}

func openDB() (*bolt.DB, error) {
	db, err := bolt.Open("urlshort.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDB(db *bolt.DB) error {
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("url_mappings"))
		if err != nil {
			return err
		}
		if err = b.Put([]byte("/bbolt"), []byte("https://github.com/coreos/bbolt")); err != nil {
			return err
		}
		if err = b.Put([]byte("/bbolt-issues"), []byte("https://github.com/coreos/bbolt/issues")); err != nil {
			return err
		}
		return nil
	})
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
