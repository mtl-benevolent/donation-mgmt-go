#!/bin/bash
set -e

echo "[INFO] Waiting for DB"
while !</dev/tcp/$DB_HOST/$DB_PORT; do sleep 2; done;

# Leaving some time for the init script to be applied
sleep 10

echo "[INFO] Migrating the DB to the latest version"
npx prisma migrate deploy

echo "[INFO] Seeding the DB"
make seed
