A command line utility that downsamples the rows of the input file or Unix pipe.
It works like Unix `cat` but dropping random rows, or like `head` or `tail` but
taking random rows.

To build it run the `make` command that needs [Go](https://go.dev/) installed.
The `make test-samples` tests additionally use [R](https://cran.r-project.org/).

**Usage example:**

```shell
$ rand -l main.go
    50          }
    52          // seen > c.size items, randomly replace with new ones
    54                  i := rand.Intn(c.size)
    66  type Args struct {
    67          file    *os.File
    88                  fmt.Printf("\nExamples:\n")
   100                  exit(fmt.Errorf("fraction of rows needs to be a value between 0 and 100 (%%), got %v", frac))
   141                  return
   144          // using number of lines option
   150
$ cat main.go | rand -p 0.1   
import (
        "sort"
type Printer struct {

        } else {
func newLineCache(size int) LineCache {
                return
        if rand.Float64() < (float64(c.size) / c.counter) {
        return c.cache

}
        fmt.Fprintln(os.Stderr, msg)      
```
