FROM golang:1.18-alpine as build

ARG BIN_FILE
ARG APP_FOLDER

ENV BIN_FILE $BIN_FILE
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ./cmd/${APP_FOLDER}

# На выходе тонкий образ
FROM alpine:latest

ARG CONFIG_FILE_EXTERNAL

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ARG BIN_FILE
ENV BIN_FILE $BIN_FILE
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/app/config.toml
COPY ${CONFIG_FILE_EXTERNAL} ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
