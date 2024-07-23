FROM golang:1.22-alpine as builder

RUN apk add --no-cache git make

WORKDIR /notifier

COPY notifier/go.mod notifier/go.sum ./
RUN go mod download

COPY notifier/ ./

RUN make build

FROM alpine
WORKDIR /notifier
COPY --from=builder /notifier/bin/notifier ./app

ENTRYPOINT ["/notifier/app"]
