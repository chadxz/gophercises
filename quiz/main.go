package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	questionsFileNamePtr := flag.String("input", "problems.csv",
		"File to read questions and answers from")

	flag.Parse()

	log.Print("The file name is ", *questionsFileNamePtr)
	questionsFile, err := os.Open(*questionsFileNamePtr)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(questionsFile)
	in := bufio.NewReader(os.Stdin)
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		fmt.Print(rec[0], "? ")
		answer, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		answer = strings.TrimSpace(answer)
		if answer != rec[1] {
			fmt.Printf("nope. %s != %s\n", answer, rec[1])
		}
	}
}
