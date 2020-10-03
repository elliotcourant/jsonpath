.PHONY: default tests dependencies

default: tests

dependencies:
	go get -t -v ./...

tests:
	go test -v ./...

coverage:
	./coverage.sh