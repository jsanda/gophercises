package camelcase

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type camelCaseReader struct {}

func NewCamelCaseReader() (*camelCaseReader, error) {
	return &camelCaseReader{}, nil
}

func (c *camelCaseReader) Run() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("input: ")
		text, err := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if err != nil {
			return err
		}
		fmt.Printf("%s has %d words\n\n", text, getWordCount(text))
	}
}

func getWordCount(s string) int {
	count := 1
	for i :=0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' && i > 0{
			count++
		}
	}
	return count
}
