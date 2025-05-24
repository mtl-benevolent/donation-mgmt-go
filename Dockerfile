FROM golang:1.24.2-bookworm AS deps

ENV NODE_MAJOR=22

RUN apt-get update && \
  apt-get install git ca-certificates make curl gnupg && \
  update-ca-certificates && \
  mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
  echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
  apt-get update -y && \
  apt-get install nodejs -y && \
  go install github.com/go-delve/delve/cmd/dlv@v1.24.2
  
WORKDIR /build
RUN git init --quiet

COPY go.mod go.sum package*.json Makefile ./

RUN go mod download && \
  npm ci && \
  make deps

FROM deps AS build

WORKDIR /build
ENV CGO_ENABLED=1

COPY ./ ./

RUN make build

# This CMD is used for local development
CMD [ "sh", "-c", "make build && /go/bin/dlv exec ./dist/api --headless --api-version 2 --continue --accept-multiclient --listen \"0.0.0.0:18000\"" ]

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=build /build/dist/api ./api

EXPOSE 8000

CMD [ "/app/api" ]
