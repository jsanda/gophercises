package exercises

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

type quiz struct {}

func NewQuiz() *quiz {
	return &quiz{}
}

func (*quiz) Run() error {
	fileName := "exercises/problems.csv"
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", fileName, err)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	inputReader := bufio.NewReader(os.Stdin)
	var right, wrong int32

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")
		if len(tokens) != 2 {
			return fmt.Errorf("Parse error at line: %s", line)
		}
		question := tokens[0]
		answer := tokens[1]

		fmt.Printf("%s? ", question)

		input, err := inputReader.ReadString('\n')

		if err != nil {
			return fmt.Errorf("Failed to read user input: %s", err)
		}

		if strings.TrimRight(input, "\n") == answer {
			right++
		} else {
			wrong++
		}
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("Reading %s failed: %s", fileName, err)
	}

	fmt.Printf("\nright: %d, wrong: %d\n", right, wrong)

	return nil
}