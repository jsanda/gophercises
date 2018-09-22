package main

import (
	"github.com/jsanda/gophercises/camelcase"
	"github.com/jsanda/gophercises/exercises"
	"flag"
	"fmt"
	"github.com/jsanda/gophercises/quiz"
	"github.com/jsanda/gophercises/sitemap"
	"github.com/jsanda/gophercises/urlshort"
	"github.com/jsanda/gophercises/cyoa"
	"github.com/jsanda/gophercises/linkparser"
	"log"
)

type noOp struct {}

func (*noOp) Run() error {
	fmt.Println("No exercise chosen")
	return nil
}

func main() {
	opt := flag.String("exercise", "", "Specifies the exercise to run")
	flag.Parse()

	exercise, err := getExercise(opt)
	if err != nil {
		log.Fatalf("Failed to get exercise: %s", err)
	}
	if err:= exercise.Run(); err != nil {
		log.Fatalf("%s failed: %s", *opt, err)
	}
}

func getExercise(exercise *string) (exercises.Exercise, error) {
	switch *exercise {
	case "quiz":
		return quiz.NewQuiz()
	case "urlshort":
		return urlshort.NewUrlShortener()
	case "cyoa":
		return cyoa.NewAdventure()
	case "linkparser":
		return linkparser.NewParser()
	case "sitemap":
		return sitemap.NewBuilder()
	case "camelcase":
		return camelcase.NewCamelCaseReader()
	default:
		return &noOp{}, nil
	}
}
