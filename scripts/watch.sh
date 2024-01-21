#!/bin/bash
set -e

echo "[INFO] Waiting for DB"
while !</dev/tcp/$DB_HOST/$DB_PORT; do sleep 1; done;

echo "[INFO] Migrating the DB to the latest version"
npx prisma migrate deploy

/go/bin/air -c .air.toml
