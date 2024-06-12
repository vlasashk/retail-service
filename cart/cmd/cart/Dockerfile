FROM golang:1.22-alpine as builder

RUN apk add --no-cache git make protoc

WORKDIR /cart

COPY cart/go.mod cart/go.sum ./
RUN go mod download

COPY cart/ ./
COPY proto_contracts /proto_contracts

RUN make build

FROM alpine
WORKDIR /cart
COPY --from=builder /cart/bin/cart ./app

EXPOSE 8082
ENTRYPOINT ["/cart/app"]
