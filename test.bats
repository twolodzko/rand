#!/usr/bin/env bats

@test "Return 10 lines by default" {
    [ "$( ./rand main.go | wc -l |  cut -f1 )" -eq 10 ]
}

@test "Return custom number of lines with -n 15" {
    [ "$( ./rand -n 15 main.go | wc -l | cut -f1 )" -eq 15 ]
}

@test "With -p 0.9 return at least 10 lines" {
    [ "$( ./rand -p 0.9 main.go | wc -l | cut -f1 )" -ge 10 ]
}

@test "Return all the lines for a small file" {
    [ "$( ./rand go.mod | wc -l | cut -f1 )" -eq "$( cat go.mod | wc -l | cut -f1 )" ]
}

@test "Return all the lines for a small file for -p 1" {
    [ "$( ./rand -p 1 go.mod | wc -l | cut -f1 )" -eq "$( cat go.mod | wc -l | cut -f1 )" ]
}

@test "Show line numnbers when using -l flag" {
    [ "$( ./rand -l main.go | grep -c -E "\s+[0-9]+" )" -eq 10 ]
}

@test "Returns unique lines" {
    [ "$( cat -n main.go | ./rand -n 20 | uniq | wc -l )" -eq 20 ]
}

@test "Return unique lines for -p 0.5" {
    result="$( cat -n main.go | ./rand -p 0.5 )"
    [ "$( echo "$result" | uniq | wc -l )" -eq "$( echo "$result" | wc -l )" ]
}

@test "Fixing random seed leads to the same results" {
    [ "$( ./rand -r 42 -n 10 main.go )" == "$( ./rand -r 42 -n 10 main.go )" ]
}

@test "Fixing random seed leads to the same results with -p 0.5" {
    [ "$( ./rand -r 42 -p 0.5 main.go )" == "$( ./rand -r 42 -p 0.5 main.go )" ]
}

