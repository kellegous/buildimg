bin/buildimg: main.go $(shell find pkg -type f)
	go build -o $@

clean:
	rm -rf bin