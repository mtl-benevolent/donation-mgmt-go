FROM golang:1.21.4-bookworm AS deps

ENV NODE_MAJOR=20

RUN apt-get update && \
  apt-get install git ca-certificates make curl gnupg && \
  update-ca-certificates && \
  mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
  echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
  apt-get update && \
  apt-get install nodejs -y

WORKDIR /build

COPY go.mod go.sum package*.json ./

RUN git init --quiet && \
  go mod download && \
  npm ci

FROM deps AS build

WORKDIR /build
ENV CGO_ENABLED=1

COPY ./ ./

RUN make build

FROM alpine:3.18.3

WORKDIR /app

COPY --from=build /build/dist/api ./api

EXPOSE 8000

ENTRYPOINT [ "/app/api" ]
