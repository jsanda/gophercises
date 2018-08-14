package main

import (
	"github.com/jsanda/gophercises/exercises"
	"flag"
	"fmt"
	"github.com/jsanda/gophercises/quiz"
	"github.com/jsanda/gophercises/urlshort"
	"github.com/jsanda/gophercises/cyoa"
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
		return quiz.NewQuiz()
	case "urlshort":
		return urlshort.NewUrlShortener()
	case "cyoa":
		return cyoa.NewAdventure()
	default:
		return &noOp{}
	}
}
