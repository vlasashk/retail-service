FROM golang:1.22-alpine as builder

RUN apk add --no-cache git make protoc

WORKDIR /loms

COPY loms/go.mod loms/go.sum ./
RUN go mod download && go install github.com/grpc-ecosystem/grpc-health-probe@latest

COPY loms/ ./
COPY proto_contracts /proto_contracts

RUN make build

FROM alpine
WORKDIR /loms
COPY --from=builder /loms/bin/loms ./app
COPY --from=builder /go/bin/grpc-health-probe /bin/grpc_health_probe

EXPOSE 8082
ENTRYPOINT ["/loms/app"]
