#!/bin/bash
set -e

go install github.com/cosmtrek/air@v1.27.8
go install github.com/go-delve/delve/cmd/dlv@latest
make goose

# apk add postgresql-client
echo "Waiting for DB"
while !</dev/tcp/$DB_HOST/$DB_PORT; do sleep 1; done;

echo "Migrating the DB to the latest version"
make db_up

air -c .air.toml
