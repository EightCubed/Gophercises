package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"math/rand"
)

type Quiz struct {
	csvPath        string
	questionList   [][]string
	timeLimit      time.Duration
	correctAnswers int
	randomize      bool
}

const DEFAULT_CSV_PATH = "./problems.csv"

func (q *Quiz) readCsvFile() {
	filePath := q.csvPath
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	if q.randomize {
		rand.Shuffle(len(records), func(i, j int) { records[i], records[j] = records[j], records[i] })
	}
	q.questionList = records
}

func (q *Quiz) validateQuestions() {
	for _, question := range q.questionList {
		_, err := strconv.Atoi(question[1])
		if err != nil {
			log.Fatal("Invalid question format entered")
		}
	}
}

func (q *Quiz) processInput(answer string, question []string) {
	if answer == "" {
		return
	}

	userAnswer, err := strconv.Atoi(strings.TrimSpace(answer))
	if err != nil {
		return
	}

	expectedAnswer, err := strconv.Atoi(question[1])
	if err != nil {
		log.Fatal("Error! Something went wrong")
	}

	if expectedAnswer == userAnswer {
		q.correctAnswers++
	}
}

func (q *Quiz) calculateScore() {
	fmt.Printf("You answered %d correctly out of a total %d\n", q.correctAnswers, len(q.questionList))
}

func main() {
	flagPtr := flag.String("file", DEFAULT_CSV_PATH, "a csv file in the format of 'question,answer' [ Default : problems.csv ]")
	timePtr := flag.Int("time", 30, "the time limit for the quiz in seconds [ Default : 30 ]")
	randomizePtr := flag.Bool("random", false, "shuffle the order in which questions are asked from the file [ Default : false ]")
	flag.Parse()

	if *timePtr < 0 {
		log.Fatal("Invalid time entered")
	}

	q := &Quiz{
		csvPath:   *flagPtr,
		timeLimit: time.Duration(*timePtr),
		randomize: *randomizePtr,
	}

	q.readCsvFile()
	q.validateQuestions()

	fmt.Println("Quiz has started with a time limit")
	timer := time.After(q.timeLimit * time.Second)
	for idx, question := range q.questionList {
		fmt.Printf("\t%d ) %s = ", idx+1, question[0])

		answerCh := make(chan string)
		go func() {
			var input string
			fmt.Scanln(&input)
			answerCh <- input
		}()

		select {
		case <-timer:
			fmt.Println("\nTime's up!")
			q.calculateScore()
		case answer := <-answerCh:
			q.processInput(answer, question)
		}
	}

	q.calculateScore()
}
