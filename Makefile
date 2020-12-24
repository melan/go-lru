.PHONY: test

all: test

test:
	go test -timeout 30s ./...