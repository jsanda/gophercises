package linkparser

import (
	"flag"
	"fmt"
	"bufio"
	"io"
	"os"
	"golang.org/x/net/html"
	"strings"
	"unicode"
)

type parser struct {
	Html io.Reader
	Links []*Link
}

type Link struct {
	Href string
	Text string
}

var (
	FileName = flag.String("html", "", "Specifies the HTML file to parse")
)

func NewParser() (*parser, error) {
	file, err := os.Open(*FileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %s", *FileName, err)
	}
	return NewParserForReader(bufio.NewReader(file))
}

func NewParserForReader(reader io.Reader) (*parser, error) {
	return &parser{Html: reader, Links: make([]*Link, 0)}, nil
}

func (p *parser) Run() error {
	z := html.NewTokenizer(p.Html)
	for {
		if z.Next() == html.ErrorToken {
			return z.Err()
		}
		token := z.Token()
		if (token.Type == html.CommentToken) {
			continue
		}
		if token.Type == html.StartTagToken && token.Data == "a"{
			if link, err := parseAnchor(z, token); err != nil {
				return err
			} else {
				p.Links = append(p.Links, link)
			}
		} else if token.Type == html.EndTagToken && token.Data == "html" {
			break
		}
	}

	return nil
}

func parseAnchor(z *html.Tokenizer, t html.Token) (*Link, error) {
	link := Link{}
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
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
			if !isSpace(token.Data) {
				text = append(text, strings.TrimSpace(token.Data))
			} 
		}
	}
}

func isSpace(s string) bool {
	for _, r := range []rune(s) {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
