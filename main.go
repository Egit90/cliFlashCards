package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

type AnswerInput interface {
	GetAnswer() string
}

type problem struct {
	question string
	answer   string
}

type Quiz struct {
	problems    []problem
	timeLimit   time.Duration
	timer       *time.Timer
	answerInput AnswerInput
}

type CommandLineInput struct{}

func (cli CommandLineInput) GetAnswer() string {
	var answer string
	fmt.Scanf("%s\n", &answer)
	return answer
}

func NewQuiz(csvFilename string, timeLimit int, answerInput AnswerInput) (*Quiz, error) {
	file, err := os.Open(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the CSV file %s\n", csvFilename)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("Failed to parse CSV file")
	}

	Shuffle(lines)
	problems := parseLines(lines)

	return &Quiz{
		problems:    problems,
		timeLimit:   time.Duration(timeLimit) * time.Second,
		timer:       time.NewTimer(time.Duration(timeLimit) * time.Second),
		answerInput: answerInput,
	}, nil
}

func (q *Quiz) Run() {
	counter := 0
	for i, p := range q.problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		answerCh := make(chan string)

		go func() {
			answer := q.answerInput.GetAnswer()
			answerCh <- answer
		}()

		select {
		case <-q.timer.C:
			fmt.Printf("You Scored %d out of %d\n", counter, len(q.problems))
			return
		case answer := <-answerCh:
			if answer == p.answer {
				fmt.Println("Correct!")
				counter++
			}
		}
	}
	fmt.Printf("You Scored %d out of %d\n", counter, len(q.problems))
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in format of 'problem,answer' ")
	timeLimit := flag.Int("limit", 5, "Provide time limit for quiz in seconds")
	flag.Parse()

	quiz, err := NewQuiz(*csvFilename, *timeLimit, CommandLineInput{})
	if err != nil {
		exit(err.Error())
	}

	quiz.Run()
}

func Shuffle(slice interface{}) {
	v := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)

	length := v.Len()

	rand.Shuffle(length, func(i, j int) {
		swap(i, j)
	})
}
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
