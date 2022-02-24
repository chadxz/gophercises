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

type resultset struct {
	correct   int
	incorrect int
}

func main() {
	questionsFileNamePtr := flag.String("input", "problems.csv",
		"File to read questions and answers from")

	flag.Parse()

	log.Print("Reading quiz from file: ", *questionsFileNamePtr)
	records, err := readQuiz(*questionsFileNamePtr)
	check(err)

	stdin := bufio.NewReader(os.Stdin)
	results := resultset{}
	for _, rec := range records {
		fmt.Print(rec[0], "? ")
		answer, err := readAnswer(stdin)
		check(err)

		if answer == rec[1] {
			results.correct++
		} else {
			results.incorrect++
		}
	}

	fmt.Printf("%d/%d correct", results.correct,
		results.correct+results.incorrect)
}

func readQuiz(file string) ([][]string, error) {
	questionsFile, err := os.Open(file)
	if err != nil {
		return [][]string{}, err
	}
	defer closeFile(questionsFile)

	r := csv.NewReader(questionsFile)
	records, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
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
