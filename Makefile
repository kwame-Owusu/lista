BINARY_NAME=lista

GO=go
GOFMT=gofmt

build:
	$(GO) build -o bin/$(BINARY_NAME) main.go 

run: build
	./bin/$(BINARY_NAME)

all: test build

.PHONY: test
test:
	$(GO) test -v -race -buildvcs ./...

tidy:
	$(GO) mod tidy
	$(GO) fmt ./...

.PHONY: clean
clean:
	rm -rf bin/
