package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvFile := flag.String("file", "problems.csv", "Choose CSV file(format `question,answer`) to read")
	timeLimit := flag.Int("limit", 30, "Number of seconds available to take the quiz in")
	flag.Parse()
	openedFile, err := os.Open(*csvFile)
	if err != nil {
		fmt.Printf("Filed to open file : %s\n", *csvFile)
		os.Exit(1)
	}
	defer openedFile.Close()
	reader := csv.NewReader(openedFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("error reading records: ", err)
	}
	timerExpire := make(chan bool)
	go func() {
		time.Sleep(time.Duration(*timeLimit) * time.Second)
		timerExpire <- true
	}()
	var correctAns int
	// var start string
	// fmt.Println("Press any letter to start the quiz:")
	// fmt.Scan(&start)

	problems := parseLines(records)
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)
		ansChan := make(chan string)
		go func () {
			var answer string
			fmt.Scanf("%s\n", &answer)
			ansChan <- answer
		}()
		select {
		case <-timerExpire:
			fmt.Printf("\nYou got %d/%d questions correct!\n", correctAns, len(records))
			return
		case answer := <-ansChan:
			if answer == p.answer {
				correctAns++
			} 
		}
	}
	fmt.Printf("You got %d/%d questions correct!\n", correctAns, len(records))
}

func parseLines(lines [][]string) []problem {
	p := make([]problem, len(lines))
	for i, line := range lines {
		p[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return p
}
