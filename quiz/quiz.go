package quiz

import (
	"os"
	"bufio"
	"fmt"
	"encoding/csv"
	"strings"
	"flag"
	"time"
)

var (
	fileName = flag.String("quiz", "exercises/problems.csv",
		"Specifies the quiz which should be a CSV file")
	seconds = flag.Int("time", 60,
		"Specifies the time in seconds allowed to complete the quiz")
)

type quiz struct {}

func NewQuiz() (*quiz, error) {
	return &quiz{}, nil
}

func (*quiz) Run() error {
	quiz, err := loadQuiz(fileName)
	if err != nil {
		return err
	}
	var right, wrong int32

	fmt.Printf("You have %d seconds to complete the quiz. Go!\n", *seconds)
	timer := time.NewTimer(time.Duration(*seconds) * time.Second)
	done := make(chan struct{})
	errors := make(chan error)

	go takeQuiz(quiz, &right, &wrong, done, errors)

	select {
	case <-timer.C:
		fmt.Println("\nTime is up!")
		printResults(right, wrong, int32(len(quiz)))
	case <-done:
		fmt.Println("Finished!")
		printResults(right, wrong, int32(len(quiz)))
	case err := <- errors:
		return err
	}

	return nil
}

func takeQuiz(quiz [][]string, right *int32, wrong *int32, done chan<- struct{}, errors chan<- error) {
	inputReader := bufio.NewReader(os.Stdin)
	for _, row := range quiz  {
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

func loadQuiz(fileName *string) ([][]string, error) {
	file, err := os.Open(*fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s: %s", *fileName, err)
	}
	reader := bufio.NewReader(file)
	r := csv.NewReader(reader)
	r.FieldsPerRecord = 2

	return r.ReadAll()
}

func printResults(right, wrong, numQuestions int32) {
	fmt.Printf("\nright: %d, wrong: %d, total questions: %d\n", right, wrong, numQuestions)
}