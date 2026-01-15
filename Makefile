build:
	go build -o previewer ./cmd/server

run: build
	./previewer

test:
	go test -v ./...

lint:
	golangci-lint run --config .golangci.yml ./...
