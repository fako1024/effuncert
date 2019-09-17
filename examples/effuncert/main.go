package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/fako1024/effuncert"
)

var (
	confidence float64
)

func main() {

	// Initialize logger
	l := log.New(os.Stdout, "", 0)

	// Parse flags
	flag.Float64Var(&confidence, "confidence", effuncert.OneSigma, "Estimation confidence (standard deviation equivalent interval)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		l.Fatalf("Invalid number of arguments provided. Usage: %s <nSuccess> <nTrial>", os.Args[0])
	}

	nSuccess, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		l.Fatalf("Invalid number of successes provided (integer required)")
	}
	nTrial, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		l.Fatalf("Invalid number of trials provided (integer required)")
	}

	e := effuncert.New(nSuccess, nTrial,
		effuncert.WithConfidence(confidence),
	)

	l.Printf("%s", e)
}
