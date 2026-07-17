.PHONY: run build fmt test vet clean

run:
	go run ./cmd/golethe

build:
	go build -o bin/golethe ./cmd/golethe

fmt:
	gofmt -w cmd/golethe/*.go internal/engine/*.go internal/terminal/*.go

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin coverage.out
