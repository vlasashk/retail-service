FROM golang:1.22-alpine as builder
WORKDIR /cart
COPY ../../go.sum ../../go.mod ./
RUN go mod download
COPY ../.. /cart
RUN go build -o app ./cmd/cart/main.go

FROM alpine
WORKDIR /cart
COPY --from=builder /cart/app ./

EXPOSE 8082
ENTRYPOINT ["/cart/app"]
