package sitemap

import (
	"flag"
	"fmt"
	"github.com/jsanda/gophercises/linkparser"
	"log"
	"net/http"
	"net/url"
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
	uris.enqueue(b.initalURL.RequestURI())
	b.scrape(&uris)

	return nil
}

func (b *builder) scrape(uris *queue) error {
	fmt.Printf("URIs: %v\n\n", uris)
	if uris.isEmpty() {
		return nil
	}

	uri := uris.dequeue()
	u, err := b.resolveURL(&uri)
	if err != nil {
		return err
	}

	fmt.Printf("url = %s\n", u)

	res, err := http.Get(u.String())
	if err != nil {
		log.Printf("Failed to get %s\n", u)
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

	b.scrape(uris)

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
		fmt.Printf("enqueing %s\n", s)
		q.elements = append(q.elements, s)
		q.index[s] = true
	}
}

func (q *queue) dequeue() string {
	s := q.elements[0]
	q.elements = q.elements[1:]
	return s
}