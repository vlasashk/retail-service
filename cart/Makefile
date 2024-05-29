.PHONY: build, lint
build:
	go build -o ./bin/app ./cmd/cart/main.go

lint:
	golangci-lint run -c .golangci.yaml ./...
