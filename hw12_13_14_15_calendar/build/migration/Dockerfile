FROM golang:1.18-alpine

WORKDIR "/migrations"

COPY ./migrations .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENTRYPOINT ["/go/bin/goose", "up"]
