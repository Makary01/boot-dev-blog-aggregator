#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/sql/schema/"
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up

