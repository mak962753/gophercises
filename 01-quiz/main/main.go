package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Answer struct {
	val     string
	err     error
	timeout bool
}

func cleanAnswer(s string) string {
	return strings.ToLower(strings.Trim(s, " "))
}

func ask(question string, ch chan Answer) {
	fmt.Printf("Question: %s \n", question)
	a := ""
	_, err := fmt.Scanf("%s", &a)
	ch <- Answer{cleanAnswer(a), err, false}
}

func main() {
	const TimeLimitSec = 30

	shuffle := flag.Bool("shuffle", false, "shows if to shuffle questions")
	timelimit := flag.Int("timelimit", TimeLimitSec, "time limit for the test in seconds (30 by default)")

	flag.Parse()

	f, err := os.Open("../tests.csv")
	defer func() { _ = f.Close() }()

	if err != nil || f == nil {
		fmt.Println("Failed to read questions file")
		os.Exit(1)
	}

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println("Failed to read questions file")
		return
	}

	total, questions := 0, make([][]string, 0, len(rows))

	for _, q := range rows {
		if len(q) < 2 {
			continue
		}
		q[1] = cleanAnswer(q[1])
		questions = append(questions, q)
		total++
	}

	if *shuffle {
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}

	correct, wrong := 0, total
	answers := make(chan Answer)

	fmt.Printf("\nPress a key to start a test (you have %d seconds to pass it).\n", *timelimit)
	_, _ = fmt.Scanln()

	time.AfterFunc(time.Second*TimeLimitSec, func() {
		answers <- Answer{
			val:     "",
			err:     nil,
			timeout: true,
		}
	})

	for _, q := range rows {
		go ask(q[0], answers)
		a := <-answers
		if a.timeout {
			break
		}
		if a.err == nil && a.val == q[1] {
			correct++
			wrong--
		}
	}
	_, _ = fmt.Printf("\nResults\nCorrect: %d\nIncorrect: %d\nTotal: %d\n", correct, wrong, total)

}
