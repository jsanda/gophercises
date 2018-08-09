package main

import (
	"github.com/jsanda/gophercises/exercises"
	"flag"
	"fmt"
)

type noOp struct {}

func (*noOp) Run() error {
	fmt.Println("No exercise chosen")
	return nil
}

func main() {
	opt := flag.String("exercise", "", "Specifies the exercise to run")
	flag.Parse()

	exercise := getExercise(opt)
	exercise.Run()
}

func getExercise(exercise *string) exercises.Exercise {
	switch *exercise {
	case "quiz":
		return exercises.NewQuiz()
	case "urlshort":
		return exercises.NewUrlShortener()
	default:
		return &noOp{}
	}
}
