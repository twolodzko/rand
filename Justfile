build:
	go build -o rand main.go

test: build
	bats test.bats

test-samples: build
	#!/bin/bash
	set -eu -o pipefail

	# setup
	rm -rf result.data
	rm -rf example.data
	for i in {1..100}; do \
		echo $i >> example.data; \
	done

	echo "Test: rand -p 0.1 returns uniform results"
	time for _ in {1..5000}; do \
		./rand -p 0.1 example.data >> result.data; \
	done
	R --vanilla -q -f test-uniformity.R
	rm -rf result.data

	echo "Test: rand -n 10 returns uniform results"
	time for _ in {1..5000}; do \
		./rand -n 10 example.data >> result.data; \
	done
	R --vanilla -q -f test-uniformity.R
	rm -rf result.data

	# cleanup
	rm -rf example.data

	echo "OK"
