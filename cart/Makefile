.PHONY: build
build:
	go build -o ./bin/cart ./cmd/cart/main.go

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yaml ./...

.PHONY: unit-test
unit-test:
	go test -race -coverprofile ./coverage/cover.out ./internal/... && \
    go tool cover -html=./coverage/cover.out -o ./coverage/cover.html && \
    open ./coverage/cover.html && \
    rm ./coverage/cover.out

.PHONY: integration-test
#https://github.com/golang/go/issues/65653
integration-test: export GOEXPERIMENT=nocoverageredesign
integration-test:
	go test -race -tags=integration -coverprofile ./coverage/cover.out -coverpkg ./internal/... ./... && \
	go tool cover -html=./coverage/cover.out -o ./coverage/cover.html && \
	open ./coverage/cover.html && \
	rm ./coverage/cover.out
