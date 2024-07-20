.PHONY: build
build:
	go build -o ./bin/notifier ./cmd/notifier/main.go

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml ./...
