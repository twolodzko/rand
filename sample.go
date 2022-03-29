package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Args struct {
	file   *os.File
	frac   float64
	seed   int64
	number bool
}

func main() {
	var (
		err    error
		rownum int = 1
	)
	args := parseArgs()

	rand.Seed(args.seed)

	scanner := bufio.NewScanner(bufio.NewReader(args.file))
	for scanner.Scan() {
		if rand.Float64() < args.frac {
			line := scanner.Text()
			if args.number {
				_, err = fmt.Fprintf(os.Stdout, "%6d\t%s\n", rownum, line)
			} else {
				_, err = fmt.Fprintln(os.Stdout, line)
			}
			if err != nil {
				exit(err)
			}
		}
		rownum++
	}
}

func parseArgs() Args {
	var (
		file   *os.File
		frac   float64
		seed   int64
		number bool
		err    error
	)

	flag.Usage = func() {
		fmt.Println("Sample fraction of rows of the input")
		fmt.Printf("\nUsage:\n  %s [OPTIONS]... [FILE]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Float64Var(&frac, "p", 10, "percentage of rows to keep")
	flag.Int64Var(&seed, "r", time.Now().UnixNano(), "random seed, unix time be default")
	flag.BoolVar(&number, "n", false, "number the output lines")
	flag.Parse()

	if frac < 0 || frac > 100 {
		fmt.Printf("fraction of rows needs to be a value between 0 and 100 (%%), got %v\n", frac)
		os.Exit(1)
	}

	if flag.NArg() > 0 {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			exit(err)
		}
	} else {
		file = os.Stdin
	}

	return Args{file, frac / 100, seed, number}
}

func exit(msg error) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
