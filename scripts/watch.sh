#!/bin/bash
set -e

go install github.com/cosmtrek/air@latest

air -c .air.toml
