SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

rand: main.go
	go build -o rand main.go

.PHONY: test
test: rand
	@ time for _ in {1..100}; do \
		bash test.sh; \
	done

.PHONY: test-samples
test-samples: rand
	@ rm -rf example.data
	@ for i in {1..100}; do \
		echo $$i >> example.data; \
	done

	time for _ in {1..50000}; do \
		./rand -p 0.1 example.data >> result.data; \
	done
	R --vanilla -q -f test-uniformity.R
	rm -rf result.data

	time for _ in {1..50000}; do \
		./rand -n 10 example.data >> result.data; \
	done
	R --vanilla -q -f test-uniformity.R
	rm -rf result.data

	@ rm -rf example.data
