.PHONY: all clean

all: bin/gz

bin/%: cmd/%/main.go $(shell find internal -type f -name '*.go')
	go build -o $@ ./cmd/$*

clean:
	rm -rf bin