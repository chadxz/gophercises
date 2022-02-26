package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func main() {
	inputFileNamePtr := flag.String("input", "problems.csv",
		"File to read questions and answers from")

	flag.Parse()

	log.Print("Reading quiz from file: ", *inputFileNamePtr)
	problems, err := readQuiz(*inputFileNamePtr)
	check(err)

	stdin := bufio.NewReader(os.Stdin)

	correct := 0
	for _, problem := range problems {
		fmt.Printf("%s? ", problem.question)
		answer, err := readAnswer(stdin)
		check(err)

		if answer == problem.answer {
			correct++
		}
	}

	fmt.Printf("%d/%d correct\n", correct, len(problems))
}

func readQuiz(file string) ([]problem, error) {
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
