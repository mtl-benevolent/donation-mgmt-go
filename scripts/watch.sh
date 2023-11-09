#!/bin/bash
set -e

go install github.com/cosmtrek/air@v1.27.8
go install github.com/go-delve/delve/cmd/dlv@latest

echo "[INFO] Waiting for DB"
while !</dev/tcp/$DB_HOST/$DB_PORT; do sleep 1; done;

echo "[INFO] Migrating the DB to the latest version"
npx prisma migrate deploy

air -c .air.toml
