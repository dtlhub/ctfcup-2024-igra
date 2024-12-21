#!/usr/bin/env bash

set -e

LEVEL="test"

AUTH_TOKEN=dtl:12345 go run cmd/server/main.go -s ":8085" -h -l "$LEVEL"
