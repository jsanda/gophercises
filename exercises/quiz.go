package exercises

import (
	"os"
	"bufio"
	"fmt"
	"encoding/csv"
	"io"
	"strings"
	"flag"
)

type quiz struct {}

func NewQuiz() *quiz {
	return &quiz{}
}

func (*quiz) Run() error {
	fileName := flag.String("quiz", "exercises/problems.csv",
		"Specifies the quiz which should be a CSV file")

	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", *fileName, err)
	}

	reader := bufio.NewReader(file)
	inputReader := bufio.NewReader(os.Stdin)
	var right, wrong int32

	r := csv.NewReader(reader)
	r.FieldsPerRecord = 2

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Failed to parse %s: %s", *fileName, err)
		}
		question := row[0]
		answer := row[1]

		fmt.Printf("%s? ", question)
		input, err := inputReader.ReadString('\n')

		if err != nil {
			return fmt.Errorf("Failed to read input: %s", err)
		}

		if strings.TrimRight(input, "\n") == answer {
			right++
		} else {
			wrong++
		}
	}

	fmt.Printf("\nright: %d, wrong: %d\n", right, wrong)

	return nil
}