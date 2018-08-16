package linkparser_test

import (
	"testing"
	"github.com/jsanda/gophercises/linkparser"
)

func TestSimpleLink(t *testing.T) {
	linkparser.FileName = htmlFile("testdata/ex1.html")
	parser := linkparser.NewParser()

	if err := parser.Run(); err != nil {
		t.Fatal(err)
	}

	expected := []*linkparser.Link {
		{
			Href: "/other-page",
			Text: "A link to another page",
		},
	}

	n := len(parser.Links)
	if n == 0 {
		t.Fatal("Failed to parse any links")
	} else if n > 1 {
		t.Errorf("Parsed %d links but expected %d", n, len(expected))
	} else if *parser.Links[0] != *expected[0] {
		t.Errorf("expected: %+v, actual: %+v", *expected[0], *parser.Links[0])
	}
}

func TestMultipleLinksWithNestedElements(t *testing.T) {
	linkparser.FileName = htmlFile("testdata/ex2.html")
	parser := linkparser.NewParser()

	if err := parser.Run(); err != nil {
		t.Fatal(err)
	}

	expected := []*linkparser.Link {
		{
			Href: "https://www.twitter.com/joncalhoun",
			Text: "Check me out on twitter",
		},
		{
			Href: "https://github.com/gophercises",
			Text: "Gophercises is on Github !",
		},
	}

	n := len(parser.Links)
	if n == 0 {
		t.Fatal("Failed to parse any links")
	} else if len(expected) != len(parser.Links) {
		t.Fatal("Wrong number of links")
	} else {
		compareLinks(t, expected, parser.Links)
	}
}

func TestLinksNestInOtherElements(t *testing.T) {
	linkparser.FileName = htmlFile("testdata/ex3.html")
	parser := linkparser.NewParser()

	if err := parser.Run(); err != nil {
		t.Fatal(err)
	}

	expected := []*linkparser.Link {
		{
			Href: "#",
			Text: "Login",
		},
		{
			Href: "/lost",
			Text: "Lost? Need help?",
		},
		{
			Href: "https://twitter.com/marcusolsson",
			Text: "@marcusolsson",
		},
	}

	n := len(parser.Links)
	if n == 0 {
		t.Fatal("Failed to parse any links")
	} else if len(expected) != len(parser.Links) {
		t.Fatal("Wrong number of links")
	} else {
		compareLinks(t, expected, parser.Links)
	}
}

func TestLinkWithCommentInBody(t *testing.T) {
	linkparser.FileName = htmlFile("testdata/ex4.html")
	parser := linkparser.NewParser()

	if err := parser.Run(); err != nil {
		t.Fatal(err)
	}

	expected := []*linkparser.Link {
		{
			Href: "/dog-cat",
			Text: "dog cat",
		},
	}

	n := len(parser.Links)
	if n == 0 {
		t.Fatal("Failed to parse any links")
	} else if len(expected) != len(parser.Links) {
		t.Fatal("Wrong number of links")
	} else {
		compareLinks(t, expected, parser.Links)
	}
}

func compareLinks(t *testing.T, expected []*linkparser.Link, actual []*linkparser.Link) {
	for i := range expected {
		if *expected[i] != *actual[i] {
			t.Errorf("i: %d, expected: %+v, actual: %+v", i, *expected[i], *actual[i])
		}
	}
}

func htmlFile(f string) *string {
	return &f
}
