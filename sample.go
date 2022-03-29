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
	number bool
}

func (p Printer) Print(rownum int, line string) error {
	var err error
	if p.number {
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

type LineCache struct {
	cache   []Line
	size    int
	counter float64
}

func newLineCache(size int) LineCache {
	return LineCache{[]Line{}, size, 0}
}

func (c *LineCache) Add(elem string) {
	c.counter++
	line := Line{int(c.counter), elem}

	// seen < c.size items
	if c.counter <= float64(c.size) {
		c.cache = append(c.cache, line)
		return
	}

	// seen > c.size items, randomly replace with new ones
	if rand.Float64() < (float64(c.size) / c.counter) {
		i := rand.Intn(c.size)
		c.cache[i] = line
	}
}

func (c LineCache) Lines() []Line {
	sort.Slice(c.cache, func(i, j int) bool {
		return c.cache[i].rownum < c.cache[j].rownum
	})
	return c.cache
}

type Args struct {
	file    *os.File
	frac    float64
	size    int
	seed    int64
	lineNum bool
}

func parseArgs() Args {
	var (
		file    *os.File
		frac    float64
		size    int
		seed    int64
		lineNum bool
		err     error
	)

	flag.Usage = func() {
		fmt.Println("Sample fraction of rows of the input")
		fmt.Printf("\nUsage:\n  %s [OPTIONS]... [FILE]\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  sample -l /etc/hosts\n")
		fmt.Printf("  cat /etc/hosts | sample -p 50\n")
	}

	flag.IntVar(&size, "n", 10, "number of lines to sample; ignored when -p is greater than 0")
	flag.Float64Var(&frac, "p", 0, "percentage of rows (0-100) to keep; used instead of -n when -p is greater than 0")
	flag.Int64Var(&seed, "r", time.Now().UnixNano(), "random seed, unix time be default")
	flag.BoolVar(&lineNum, "l", false, "show line numbers")
	flag.Parse()

	if frac < 0 || frac > 100 {
		exit(fmt.Errorf("fraction of rows needs to be a value between 0 and 100 (%%), got %v", frac))
	}

	if flag.NArg() > 0 {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			exit(err)
		}
	} else {
		file = os.Stdin
	}

	return Args{file, frac / 100, size, seed, lineNum}
}

func exit(msg error) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	args := parseArgs()

	rand.Seed(args.seed)

	scanner := bufio.NewScanner(bufio.NewReader(args.file))
	printer := Printer{args.lineNum}
	rownum := 1

	// using percentage option
	if args.frac > 0 {
		for scanner.Scan() {
			if rand.Float64() < args.frac {
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

	// using number of lines option
	cache := newLineCache(args.size)
	for scanner.Scan() {
		line := scanner.Text()
		cache.Add(line)
	}

	// print the collected lines
	for _, line := range cache.Lines() {
		printer.Print(line.rownum, line.value)
	}
}
