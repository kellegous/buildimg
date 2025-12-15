all: bin/buildimg

.PHONY: test clean all

bin/buildimg: main.go $(shell find internal builder -name '*.go')
	go build -o $@

test:
	go test ./...

clean:
	rm -rf bin