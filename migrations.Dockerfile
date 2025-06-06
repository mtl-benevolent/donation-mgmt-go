FROM golang:1.24.2-bookworm

ENV NODE_MAJOR=22

RUN apt-get update && \
  apt-get install git ca-certificates make curl gnupg && \
  update-ca-certificates && \
  mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
  echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
  apt-get update -y && \
  apt-get install nodejs -y

WORKDIR /build

COPY go.mod go.sum package*.json ./

RUN git init --quiet && \
  go mod download && \
  npm ci

COPY ./ ./

ENTRYPOINT ["/build/scripts/migrate-and-seed.sh"]
