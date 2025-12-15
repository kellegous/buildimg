bin/buildimg: main.go $(shell find internal builder -name '*.go')
	go build -o $@

clean:
	rm -rf bin