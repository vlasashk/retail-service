.PHONY: build, lint, coverage
build:
	go build -o ./bin/cart ./cmd/cart/main.go

lint:
	golangci-lint run -c .golangci.yaml ./...

coverage:
	go test -coverprofile=cover.out ./... && \
    go tool cover -html=cover.out -o cover.html && \
    open cover.html && \
    rm cover.out

