SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

sample: sample.go
	go build sample.go

test: sample
	@ for _ in {1..100}; do \
		bash test.sh; \
	done
