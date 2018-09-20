package sitemap

import (
	"encoding/xml"
	"flag"
	"github.com/jsanda/gophercises/linkparser"
	"net/http"
	"net/url"
	"os"
)

type builder struct {
	initalURL url.URL
}

// This is a special purpose queue for a couple reasons. First, it does not allow duplicates.
// Secondly, it does not allow an element to be added again even after it has been removed.
type queue struct {
	elements []string
	index map[string]bool
}

type sitemapUrl struct {
	Loc string 			`xml:"loc"`
}

type urlset struct {
	NS   string       	`xml:"xmlns,attr"`
	URLs []sitemapUrl 	`xml:"url"`
}

var (
	address = flag.String("url", "", "Specifies the site to parse")
)

func NewBuilder() (*builder, error) {
	if u, err := url.Parse(*address); err != nil {
		return nil, err
	} else {
		return &builder{ initalURL: *u}, nil
	}
}

func (b *builder) Run() error {
	uris := queue{
		elements: make([]string, 0),
		index: make(map[string]bool),
	}
	urlSet := urlset{
		NS:   "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs: make([]sitemapUrl, 0),
	}
	uris.enqueue(b.initalURL.RequestURI())
	b.scrape(&uris, &urlSet)

	f, err := os.Create("sitemap.xml")
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := xml.NewEncoder(f)
	encoder.Indent(" ", "  ")
	if err = encoder.Encode(urlSet); err != nil {
		return err
	}

	return nil
}

func (b *builder) scrape(uris *queue, urlSet *urlset) error {
	if uris.isEmpty() {
		return nil
	}

	uri := uris.dequeue()
	u, err := b.resolveURL(&uri)
	if err != nil {
		return err
	}
	urlSet.URLs = append(urlSet.URLs, sitemapUrl{Loc: u.String()})

	res, err := http.Get(u.String())
	if err != nil {
		return err
	}

	parser, err := linkparser.NewParserForReader(res.Body)
	if err != nil {
		return err
	}

	parser.Run()
	for _, link := range parser.Links {
		// skip any links that do not have an href attribute declared
		if link.Href == "" {
			continue
		}

		nextURL, err := b.resolveURL(&link.Href)
		if err != nil {
			continue
		}
		if nextURL.Hostname() != b.initalURL.Hostname() {
			continue
		}
		uris.enqueue(nextURL.RequestURI())
	}

	b.scrape(uris, urlSet)

	return nil
}

func (b *builder) resolveURL(href *string) (*url.URL, error) {
	u, err := url.Parse(*href)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		u, err = b.initalURL.Parse(*href)
		if err != nil {
			return nil, err
		}
		return u, nil
	} else {
		return u, nil
	}
}

func (q *queue) isEmpty() bool {
	return len(q.elements) == 0
}

func (q *queue) enqueue(s string) {
	if _, found := q.index[s]; !found {
		q.elements = append(q.elements, s)
		q.index[s] = true
	}
}

func (q *queue) dequeue() string {
	s := q.elements[0]
	q.elements = q.elements[1:]
	return s
}