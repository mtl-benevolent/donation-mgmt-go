FROM golang:1.21-alpine AS deps

RUN apk update && apk upgrade && \
  apk add --no-cache git make bash

WORKDIR /build

COPY go.mod go.sum ./

RUN git init --quiet && \
  go mod download

FROM deps AS build

WORKDIR /build

COPY ./ ./

RUN make build

FROM alpine:3.18.3

WORKDIR /app

COPY --from=build /build/dist/api ./api

EXPOSE 8000

ENTRYPOINT [ "/app/api" ]
