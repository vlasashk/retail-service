.PHONY: build, lint, coverage
build:
	go build -o ./bin/cart ./cmd/cart/main.go

lint:
	golangci-lint run -c .golangci.yaml ./...

coverage:
	go test -race -coverprofile cover.out ./internal/... && \
    go tool cover -html=cover.out -o cover.html && \
    open cover.html && \
    rm cover.out
