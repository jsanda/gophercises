package cyoa

import (
	"net/http"
	"fmt"
	"html/template"
	"log"
	"encoding/json"
	"io/ioutil"
)

type adventure struct {}

type Story struct {
	Arcs map[string]*Arc
}

type Arc struct {
	Title string		`json:"title"`
	Story []string		`json:"story"`
	Options []*Option	`json:"options"`
}

type Option struct {
	Text string			`json:"text"`
	ArcName string		`json:"arc"`
}

func NewAdventure() *adventure {
	return &adventure{}
}

func (*adventure) Run() error {
	t, err := template.New("story.html").ParseFiles("cyoa/story.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %s", err)
	}

	fmt.Println("Starting the server on :8080")
	mux := defaultMux(t)
	log.Fatal(http.ListenAndServe(":8080", mux))

	return nil
}

func defaultMux(t *template.Template) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", storyHandler(t))
	return mux
}

func storyHandler(t *template.Template) http.HandlerFunc {
	story, err := loadStory()
	if err != nil {
		log.Fatalf("Failed to load story: %s", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		arcName := r.URL.Query().Get("arc")
		if arcName == "" {
			arcName = "intro"
		}
		arc := story.Arcs[arcName]
		if arc == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			if err = t.Execute(w, arc); err != nil {
				log.Printf("There was an error rendering the template: %s", err)
			}
		}
	}
}

func loadStory() (*Story, error) {
	fileName := "cyoa/gopher.json"
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var story = Story{}
	var objmap map[string]*json.RawMessage

	if err = json.Unmarshal(b, &objmap); err != nil {
		return nil, err
	}
	story.Arcs = make(map[string]*Arc)
	for k, v := range objmap {
		var arc Arc
		if err = json.Unmarshal(*v, &arc); err != nil {
			return nil, err
		}
		story.Arcs[k] = &arc
	}
	return &story, nil
}
