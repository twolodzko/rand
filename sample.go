package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	var (
		frac   float64
		seed   int64
		number bool
		in     *os.File
		err    error
	)

	flag.Usage = func() {
		fmt.Println("Sample fraction of rows of the input")
		fmt.Printf("\nUsage:\n  %s [OPTIONS]... [FILE]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Float64Var(&frac, "f", 10, "fraction (%) of rows to sample")
	flag.Int64Var(&seed, "s", time.Now().UnixNano(), "random seed")
	flag.BoolVar(&number, "n", false, "number the output lines")
	flag.Parse()

	if frac < 0 || frac > 100 {
		fmt.Printf("fraction of rows needs to be a value between 0 and 100 (%%), got %v\n", frac)
		os.Exit(1)
	}
	rand.Seed(seed)

	if flag.NArg() > 0 {
		in, err = os.Open(flag.Arg(0))
		if err != nil {
			exit(err)
		}
	} else {
		in = os.Stdin
	}

	scanner := bufio.NewScanner(bufio.NewReader(in))
	var row int
	for scanner.Scan() {
		row++
		if rand.Float64() < frac/100 {
			line := scanner.Text()
			if number {
				_, err = fmt.Fprintf(os.Stdout, "%6d\t%s\n", row, line)
			} else {
				_, err = fmt.Fprintln(os.Stdout, line)
			}
			if err != nil {
				exit(err)
			}
		}
	}
}

func exit(msg error) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
