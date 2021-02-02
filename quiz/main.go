package main

import (
	"fmt"
	"os"
	"flag"
	"encoding/csv"
	"strings"
	"time"
)
func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file of the format (question, answeer)")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Unable to open the csv file: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit("Unable to Read the provided CSV file")
	}

	correct := 0
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		answerChannel := make(chan string)
		go func(){
			var answer string
			fmt.Scanf("%s \n", &answer)
			answerChannel <- answer
		}()
		select {
		case <- timer.C :
			fmt.Printf("\nYou scored %d out of %d \n", correct, len(problems))
			return 
		case answer := <- answerChannel:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d \n", correct, len(problems))

}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return result
}