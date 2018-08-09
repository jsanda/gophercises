package exercises

import (
	"os"
	"bufio"
	"fmt"
	"encoding/csv"
	"io"
	"strings"
	"flag"
	"time"
)

type quiz struct {}

func NewQuiz() *quiz {
	return &quiz{}
}

func (*quiz) Run() error {
	fileName := flag.String("quiz", "exercises/problems.csv",
		"Specifies the quiz which should be a CSV file")
	seconds := flag.Int("time", 60,
		"Specifies the time in seconds allowed to complete the quiz")

	flag.Parse()

	var right, wrong int32

	fmt.Printf("You have %d seconds to complete the quiz. Go!\n", *seconds)
	timer := time.NewTimer(time.Duration(*seconds) * time.Second)
	done := make(chan struct{})
	errors := make(chan error)

	go takeQuiz(fileName, &right, &wrong, done, errors)

	select {
	case <-timer.C:
		fmt.Println("\nTime is up!")
		fmt.Printf("\nright: %d, wrong: %d\n", right, wrong)
	case <-done:
		fmt.Println("Finished!")
		fmt.Printf("\nright: %d, wrong: %d\n", right, wrong)
	case err := <- errors:
		return err
	}

	return nil
}

func takeQuiz(fileName *string, right *int32, wrong *int32, done chan<- struct{}, errors chan<- error) {
	file, err := os.Open(*fileName)
	if err != nil {
		errors <- fmt.Errorf("Failed to open %s: %s", *fileName, err)
		return
	}

	reader := bufio.NewReader(file)
	inputReader := bufio.NewReader(os.Stdin)

	r := csv.NewReader(reader)
	r.FieldsPerRecord = 2

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors <- fmt.Errorf("Failed to parse %s: %s", *fileName, err)
			return
		}
		question := row[0]
		answer := row[1]

		fmt.Printf("%s? ", question)
		input, err := inputReader.ReadString('\n')

		if err != nil {
			errors <- fmt.Errorf("Failed to read input: %s", err)
			return
		}

		if strings.TrimRight(input, "\n") == answer {
			*right++
		} else {
			*wrong++
		}
	}

	done <- struct{}{}
}