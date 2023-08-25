package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	file := flag.String("file", "problems.csv", "Choose CSV file to read")
	countdown := flag.Int("time", 30, "Number of seconds available to take the quiz in")
	flag.Parse()
	openedFile, err := os.Open(*file)
	if err != nil {
		log.Fatal("can't open file: ", err)
	}
	defer openedFile.Close()
	reader := csv.NewReader(openedFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("error reading records: ", err)
	}
	timerExpire := make(chan bool)
	go func() {
		time.Sleep(time.Duration(*countdown) * time.Second)
		timerExpire <- true
	}()
	var correctAns int
	var start string
	fmt.Println("Press any letter to start the quiz:")
	fmt.Scan(&start)
OuterLoop:
	for _, record := range records {
		select {
		case <-timerExpire:
			fmt.Println("Time's up!")
			break OuterLoop
		default:
			fmt.Printf("%s: ", record[0])
			var answer string
			fmt.Scan(&answer)
			if answer == record[1] {
				correctAns++
			} else {
				continue
			}
		}
	}
	fmt.Printf("You got %d/%d questions correct!\n", correctAns, len(records))
}
