package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

type Printer struct {
	showLineNum bool
}

func (p Printer) Print(lineNum int, line string) error {
	var err error
	if p.showLineNum {
		_, err = fmt.Fprintf(os.Stdout, "%6d\t%s\n", lineNum, line)
	} else {
		_, err = fmt.Fprintln(os.Stdout, line)
	}
	return err
}

type Line struct {
	lineNum int
	value   string
}

type OnlineSampler struct {
	cache   []Line
	size    int64
	counter int64
}

func newOnlineSampler(size int64) OnlineSampler {
	return OnlineSampler{[]Line{}, size, 0}
}

// Uniformly at random add new lines to cache of size s.size
//
// See:
// https://stats.stackexchange.com/q/569647/35989
// https://en.wikipedia.org/wiki/Reservoir_sampling
func (s *OnlineSampler) Add(elem string) {
	s.counter++
	line := Line{int(s.counter), elem}

	// seen < s.size items
	if s.counter <= s.size {
		s.cache = append(s.cache, line)
		return
	}

	// seen > s.size items, randomly replace with new ones
	// we use zero indexing, so i is sampled from [0, s.counter)
	// if it falls into the [0, s.size) region we accept it as a new candidate
	// this leads to sampling with probability s.size / s.counter
	i := rand.Int63n(s.counter)
	if i < s.size {
		s.cache[i] = line
	}
}

func (s OnlineSampler) Lines() []Line {
	sort.Slice(s.cache, func(i, j int) bool {
		return s.cache[i].lineNum < s.cache[j].lineNum
	})
	return s.cache
}

type Args struct {
	prob        float64
	size        int64
	seed        int64
	showLineNum bool
	file        *os.File
}

func parseArgs() (Args, error) {
	var (
		prob        float64
		size        int64
		seed        int64
		showLineNum bool
		err         error
		file        *os.File
	)

	flag.Usage = func() {
		fmt.Println("Randomly downsample the rows of the input")
		fmt.Printf("\nUsage:\n  %s [OPTIONS]... [FILE]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  sample -l /etc/hosts\n")
		fmt.Printf("  cat /etc/hosts | sample -p 0.5\n")
	}

	flag.Int64Var(&size, "n", 10, "number of lines to sample; ignored when -p is greater than 0")
	flag.Float64Var(&prob, "p", 0, "probability of keeping each row; used instead of -n when -p is greater than 0")
	flag.Int64Var(&seed, "r", time.Now().UnixNano(), "random seed, unix time be default")
	flag.BoolVar(&showLineNum, "l", false, "show line numbers")
	flag.Parse()

	if prob < 0 || prob > 1 {
		return Args{}, fmt.Errorf("probability needs to be a value between 0 and 1, got %v", prob)
	}

	if flag.NArg() > 0 {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			return Args{}, err
		}
	} else {
		file = os.Stdin
	}

	return Args{prob, size, seed, showLineNum, file}, nil
}

func exit(msg error) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	args, err := parseArgs()
	if err != nil {
		exit(err)
	}

	rand.Seed(args.seed)

	scanner := bufio.NewScanner(bufio.NewReader(args.file))
	printer := Printer{args.showLineNum}
	lineNum := 1

	// using the percentage option
	if args.prob > 0 {
		for scanner.Scan() {
			if rand.Float64() < args.prob {
				line := scanner.Text()
				err := printer.Print(lineNum, line)
				if err != nil {
					exit(err)
				}
			}
			lineNum++
		}
		return
	}

	// using the number of lines option
	sampler := newOnlineSampler(args.size)
	for scanner.Scan() {
		line := scanner.Text()
		sampler.Add(line)
	}

	// print the collected lines
	for _, line := range sampler.Lines() {
		printer.Print(line.lineNum, line.value)
	}
}
