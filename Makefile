all: build

run:
	go run . run

build:
	go build -o bin/ .

.PHONY: clean
clean:
	rm -f db.sqlite
	rm -rf bin
	rm -rf log
	mkdir log