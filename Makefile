.PHONY: run build test clean

run:
	go run cmd/app/main.go

build:
	go build -o bin/url-shortener cmd/app/main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf bin