package linkparser

import (
	"flag"
	"fmt"
	"bufio"
	"os"
	"golang.org/x/net/html"
	"log"
	"strings"
)

type parser struct {
	Links []*Link
}

type Link struct {
	Href string
	Text string
}

var (
	fileName = flag.String("html", "", "Specifies the HTML file to parse")
)

func NewParser() *parser {
	return &parser{Links: make([]*Link, 0, 10)}
}

func (*parser) Run() error {
	file, err := os.Open(*fileName)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", *fileName, err)
	}
	reader := bufio.NewReader(file)
	z := html.NewTokenizer(reader)
	links := make([]*Link, 0)
	for {
		if z.Next() == html.ErrorToken {
			return z.Err()
		}
		token := z.Token()
		if token.Type == html.StartTagToken && token.Data == "a"{
			log.Printf("tag = %s", token.Data)
			if link, err := parseAnchor(z, token); err != nil {
				return err
			} else {
				links = append(links, link)
			}
		} else if token.Type == html.EndTagToken && token.Data == "html" {
			break
		}
	}

	fmt.Println("Links:")
	for _, link := range links {
		fmt.Printf("href: %s\ntext: %s\n\n", link.Href, link.Text)
	}
	return nil
}

func parseAnchor(z *html.Tokenizer, t html.Token) (*Link, error) {
	link := Link{}
	for _, attr := range t.Attr {
		log.Printf("attr key: %s, attr value: %s", attr.Key, attr.Val)
		if attr.Key == "href" {
			link.Href = attr.Val
			log.Printf("href: %s", attr.Val)
		}
	}
	text := make([]string, 0)
	for {
		if z.Next() == html.ErrorToken {
			return nil, z.Err()
		}
		token := z.Token()
		if token.Type == html.EndTagToken && token.Data == "a" {
			link.Text = strings.Join(text, " ")
			return &link, nil
		}
		if token.Type == html.CommentToken {
			continue
		}
		if token.Type == html.TextToken {
			text = append(text, token.Data)
		}
	}

}
