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
	lineNum bool
}

func (p Printer) Print(rownum int, line string) error {
	var err error
	if p.lineNum {
		_, err = fmt.Fprintf(os.Stdout, "%6d\t%s\n", rownum, line)
	} else {
		_, err = fmt.Fprintln(os.Stdout, line)
	}
	return err
}

type Line struct {
	rownum int
	value  string
}

type OnlineSampler struct {
	cache   []Line
	size    int
	counter float64
}

func newOnlineSampler(size int) OnlineSampler {
	return OnlineSampler{[]Line{}, size, 0}
}

// Uniformly at random add new lines to cache of size s.size
//
// See: https://stats.stackexchange.com/q/569647/35989
func (s *OnlineSampler) Add(elem string) {
	s.counter++
	line := Line{int(s.counter), elem}

	// seen < c.size items
	if s.counter <= float64(s.size) {
		s.cache = append(s.cache, line)
		return
	}

	// seen > c.size items, randomly replace with new ones
	if rand.Float64() < (float64(s.size) / s.counter) {
		i := rand.Intn(s.size)
		s.cache[i] = line
	}
}

func (s OnlineSampler) Lines() []Line {
	sort.Slice(s.cache, func(i, j int) bool {
		return s.cache[i].rownum < s.cache[j].rownum
	})
	return s.cache
}

type Args struct {
	file    *os.File
	prob    float64
	size    int
	seed    int64
	lineNum bool
}

func parseArgs() (Args, error) {
	var (
		file    *os.File
		prob    float64
		size    int
		seed    int64
		lineNum bool
		err     error
	)

	flag.Usage = func() {
		fmt.Println("Randomly downsample the rows of the input")
		fmt.Printf("\nUsage:\n  %s [OPTIONS]... [FILE]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  sample -l /etc/hosts\n")
		fmt.Printf("  cat /etc/hosts | sample -p 0.5\n")
	}

	flag.IntVar(&size, "n", 10, "number of lines to sample; ignored when -p is greater than 0")
	flag.Float64Var(&prob, "p", 0, "probability of keeping each row; used instead of -n when -p is greater than 0")
	flag.Int64Var(&seed, "r", time.Now().UnixNano(), "random seed, unix time be default")
	flag.BoolVar(&lineNum, "l", false, "show line numbers")
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

	return Args{file, prob, size, seed, lineNum}, nil
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
	printer := Printer{args.lineNum}
	rownum := 1

	// using the percentage option
	if args.prob > 0 {
		for scanner.Scan() {
			if rand.Float64() < args.prob {
				line := scanner.Text()
				err := printer.Print(rownum, line)
				if err != nil {
					exit(err)
				}
			}
			rownum++
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
		printer.Print(line.rownum, line.value)
	}
}
