#!/bin/bash
set -eu -o pipefail

if [ "$( ./rand main.go | wc -l |  cut -f1 )" -ne 10 ]; then
    echo "sample didn't return 10 lines"
    exit 1
fi

if [ "$( ./rand -n 15 main.go | wc -l | cut -f1 )" -ne 15 ]; then
    echo "sample -n 15 didn't return 15 lines"
    exit 1
fi

if [ "$( ./rand -p 0.9 main.go | wc -l | cut -f1 )" -le 10 ]; then
    echo "sample -p 0.9 returned <= 10 lines"
    exit 1
fi

if [ "$( ./rand go.mod | wc -l | cut -f1 )" -ne "$( cat go.mod | wc -l | cut -f1 )" ]; then
    echo "sample of small file didn't return all the lines"
    exit 1
fi

if [ "$( ./rand -p 1 go.mod | wc -l | cut -f1 )" -ne "$( cat go.mod | wc -l | cut -f1 )" ]; then
    echo "sample -p 1 of small file didn't return all the lines"
    exit 1
fi

if [ "$( ./rand -l main.go | grep -c -E "\s+[0-9]+" )" -ne 10 ]; then
    echo "sample -l didn't number the lines"
    exit 1
fi

if [ "$( cat -n main.go | ./rand -n 20 | uniq | wc -l )" -ne 20 ]; then
    echo "sample -n 20 didn't return unique lines"
    exit 1
fi

result="$( cat -n main.go | ./rand -p 0.5 )"
if [ "$( echo "$result" | uniq | wc -l )" -ne "$( echo "$result" | wc -l )" ]; then
    echo "sample -p 0.5 didn't return unique lines"
    exit 1
fi
unset result

if [ "$( ./rand -r 42 -n 10 main.go )" != "$( ./rand -r 42 -n 10 main.go )" ]; then
    echo "fixing random seed didn't lead to identical results for sample -n"
    exit 1
fi

if [ "$( ./rand -r 42 -p 0.5 main.go )" != "$( ./rand -r 42 -p 0.5 main.go )" ]; then
    echo "fixing random seed didn't lead to identical results for sample -p"
    exit 1
fi
