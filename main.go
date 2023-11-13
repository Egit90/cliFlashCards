package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

type AnswerInput interface {
	GetAnswer() cliResponse
}

type problem struct {
	Command     string
	Description string
	Examples    []string
}

type cliResponse struct {
	str   string
	Error error
}

type Quiz struct {
	problems    []problem
	timeLimit   time.Duration
	timer       *time.Timer
	answerInput AnswerInput
}

type CommandLineInput struct{}

func (cli CommandLineInput) GetAnswer() cliResponse {
	reader := bufio.NewReader(os.Stdin)

	answer, err := reader.ReadString('\n')

	if err != nil {
		return cliResponse{str: "", Error: err}
	}

	return cliResponse{str: strings.TrimSpace(answer), Error: nil}
}

func NewQuiz(Filename string, timeLimit int, answerInput AnswerInput) (*Quiz, error) {
	file, err := os.ReadFile(Filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open the Json file %s\n", Filename)
	}

	var problemSet []problem

	err = json.Unmarshal(file, &problemSet)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse Json file, %s", err)
	}

	Shuffle(problemSet)

	return &Quiz{
		problems:    problemSet,
		timeLimit:   time.Duration(timeLimit) * time.Second,
		timer:       time.NewTimer(time.Duration(timeLimit) * time.Second),
		answerInput: answerInput,
	}, nil
}

func (q *Quiz) Run() {
	counter := 0
	for i, p := range q.problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.Description)

		answerCh := make(chan cliResponse)

		go func() {
			Res := q.answerInput.GetAnswer()
			answerCh <- Res
		}()
		select {
		case <-q.timer.C:
			fmt.Println("Time is out 🕥")
			fmt.Printf("You Scored %d out of %d\n", counter, len(q.problems))
			return
		case answer := <-answerCh:
			if answer.Error != nil {
				println("Coundn't Read Input , ", answer.Error)
			} else if answer.str == p.Command {
				fmt.Println("Correct!✔️")
				fmt.Printf("Example Usage: %s \n", p.Examples)
				counter++
			} else {
				fmt.Printf("wrong ❌. Correct answer is %s \nExample %s \n", p.Command, strings.Join(p.Examples, " "))
			}
		}
		fmt.Println("🆕 🆕 🆕 🆕 🆕 🆕 🆕 ")
	}
	fmt.Printf("You Scored %d out of %d\n", counter, len(q.problems))
}

func main() {
	Filename := flag.String("json", "bank.json", "a json file in format of 'problem,answer' ")
	timeLimit := flag.Int("limit", 60, "Provide time limit for quiz in seconds")
	flag.Parse()

	quiz, err := NewQuiz(*Filename, *timeLimit, CommandLineInput{})
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

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
