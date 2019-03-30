BINARY_NAME=wackernews

all: build test

build:
	go build -o $(BINARY_NAME) -v

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)
