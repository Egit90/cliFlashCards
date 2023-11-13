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

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in format of 'problem,answer' ")

	timeLimit := flag.Int("limit", 5, "Provide time limit for quiz in seconds")

	flag.Parse()
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Faile to parse CSV file")
	}

	Shuffle(lines)
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	counter := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)

		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You Scored %d out of %d\n", counter, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.answer {
				fmt.Println("Correct!")
				counter++
			}

		}

	}
	fmt.Printf("You Scored %d out of %d\n", counter, len(problems))

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

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func Shuffle(slice interface{}) {
	v := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)

	length := v.Len()

	rand.Shuffle(length, func(i, j int) {
		swap(i, j)
	})
}
