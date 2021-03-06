package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	inputFileNamePtr := flag.String("input", "problems.csv",
		"File to read questions and answers from")
	timerDurationSecondsPtr := flag.Int("timer-seconds", 30,
		"Seconds to wait before the quiz ends")
	flag.Parse()

	log.Print("Reading quiz from file: ", *inputFileNamePtr)
	problems, err := readProblemsFromCSV(*inputFileNamePtr)
	check(err)

	correct := 0
	timer := time.NewTimer(
		time.Duration(*timerDurationSecondsPtr) * time.Second)
	stdin := bufio.NewReader(os.Stdin)
	for _, problem := range problems {
		answerCh := make(chan string)

		fmt.Printf("%s? ", problem.question)
		go requestAnswer(stdin, answerCh)

		select {
		case <-timer.C:
			fmt.Printf("\ntime's up! %d/%d correct\n",
				correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == problem.answer {
				correct++
			}
		}
	}

	fmt.Printf("all questions completed! %d/%d correct\n",
		correct, len(problems))
}

func requestAnswer(reader *bufio.Reader, resultCh chan<- string) {
	answer, err := readAnswer(reader)
	check(err)

	resultCh <- answer
}

func readProblemsFromCSV(file string) ([]problem, error) {
	questionsFile, err := os.Open(file)
	if err != nil {
		return []problem{}, err
	}
	defer closeFile(questionsFile)

	r := csv.NewReader(questionsFile)
	records, err := r.ReadAll()
	if err != nil {
		return []problem{}, err
	}

	return parseProblems(records), nil
}

func parseProblems(records [][]string) []problem {
	results := make([]problem, len(records))

	for i, record := range records {
		results[i] = problem{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}
	}

	return results
}

func readAnswer(reader *bufio.Reader) (string, error) {
	result, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

func closeFile(f *os.File) {
	err := f.Close()
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
