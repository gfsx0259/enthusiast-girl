FROM golang:1.16-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add gcc gettext musl-dev

COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

COPY app ./
RUN go build -o ./bin/app ./main.go

FROM alpine as runner

RUN apk --no-cache add bash git openssh-client kustomize

COPY --from=builder /usr/local/src/bin/app /

CMD ["/app"]