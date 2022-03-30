#!/bin/bash
set -eu -o pipefail

if [ "$( ./sample sample.go | wc -l |  cut -f1 -d" " )" -ne 10 ]; then
    echo "sample didn't return 10 lines"
    exit 1
fi

if [ "$( ./sample -n 15 sample.go | wc -l | cut -f1 -d" " )" -ne 15 ]; then
    echo "sample -n 15 didn't return 15 lines"
    exit 1
fi

if [ "$( ./sample -p 0.9 sample.go | wc -l | cut -f1 -d" " )" -le 10 ]; then
    echo "sample -p 0.9 returned <10 lines"
    exit 1
fi

if [ "$( ./sample go.mod | wc -l | cut -f1 -d" " )" -ne "$(wc -l go.mod | cut -f1 -d" " )" ]; then
    echo "sample of small file didn't return all the lines"
    exit 1
fi

if [ "$( ./sample -p 1 go.mod | wc -l | cut -f1 -d" " )" -ne "$(wc -l go.mod | cut -f1 -d" " )" ]; then
    echo "sample -p 1 of small file didn't return all the lines"
    exit 1
fi

if [ "$( ./sample -l sample.go | grep -c -E "\s+[0-9]+" )" -ne 10 ]; then
    echo "sample -l didn't number the lines"
    exit 1
fi

if [ "$( cat -n sample.go | ./sample -n 20 | uniq | wc -l )" -ne 20 ]; then
    echo "sample -n 20 didn't return unique lines"
    exit 1
fi

result="$( cat -n sample.go | ./sample -p 0.5 )"
if [ "$( echo "$result" | uniq | wc -l )" -ne "$( echo "$result" | wc -l )" ]; then
    echo "sample -p 0.5 didn't return unique lines"
    exit 1
fi
unset result

if [ "$( ./sample -r 42 sample.go )" != "$( ./sample -r 42 sample.go )" ]; then
    echo "fixing random seed didn't lead to identical results"
    exit 1
fi
