FROM golang:1.13.4-alpine3.10 AS build

RUN apk update && apk add --virtual build-dependencies build-base --no-cache curl \
  ca-certificates gcc autoconf automake libtool

ENV GOROOT=/usr/local/go \
  GOPATH=/app

ADD . /app/src

WORKDIR /app/src

RUN make build-binary

FROM alpine:3.10
WORKDIR /app

COPY --from=build /app/src/.env /app/
COPY --from=build /app/src/bin/server /app/

ENTRYPOINT ["/app/server"]
