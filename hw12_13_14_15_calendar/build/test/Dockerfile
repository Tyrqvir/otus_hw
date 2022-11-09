FROM golang:1.18

ARG CONFIG_FILE_EXTERNAL

WORKDIR "/tests/integration"

COPY ${CONFIG_FILE_EXTERNAL} .

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./tests/integration .

ENTRYPOINT ["go", "test", "-v"]
