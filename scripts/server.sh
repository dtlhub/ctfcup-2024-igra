#!/usr/bin/env bash

set -e

LEVEL="level-1"

AUTH_TOKEN=dtl:12345 go run cmd/server/main.go -s ":8085" -h -l "$LEVEL"
