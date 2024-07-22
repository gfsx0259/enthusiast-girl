FROM golang:1.20-alpine as modules

COPY go.mod go.sum /modules/
WORKDIR /modules

RUN go mod download

FROM golang:1.20-alpine as builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN apk --no-cache add gcc gettext musl-dev

RUN go build -o ./bin/app ./app/main.go

FROM alpine:3.20 as runner

RUN apk add bash git openssh-client docker curl envsubst openjdk11 libxml2-utils

COPY --from=builder /app/bin/app /
COPY --from=builder /app/config /config
COPY --from=builder /app/bin /bin
COPY --from=builder /app/static /static

CMD ["/app"]