package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Start starts the quiz app
func main() {

	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var input string
	var score int

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		time.Sleep(time.Duration(*limit) * time.Second)
		wg.Done()
	}()

	go func() {
		for i, line := range records {
			fmt.Printf("Problem #%d: %s = ", i+1, line[0])
			fmt.Scan(&input)

			if strings.TrimSpace(input) == line[1] {
				score++
			}
		}

		wg.Done()
	}()

	wg.Wait()
	fmt.Printf("\nYou scored %d out of %d\n", score, len(records))
}
