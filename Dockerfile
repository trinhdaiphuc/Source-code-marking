FROM golang:1.13.4-alpine3.10 AS build

RUN apk update && apk add --virtual build-dependencies build-base --no-cache \
  autoconf automake

ENV GOROOT=/usr/local/go \
  GOPATH=/app

ADD . /app/src

WORKDIR /app/src

RUN go mod download
RUN make build-binary

FROM alpine:3.21.3
WORKDIR /app

COPY --from=build /app/src/bin/server /app/

ENTRYPOINT ["/app/server"]
