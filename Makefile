bin/buildimg: main.go $(shell find internal -type f)
	go build -o $@

clean:
	rm -rf bin